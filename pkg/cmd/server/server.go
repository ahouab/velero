/*
Copyright 2017 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"

	api "github.com/heptio/ark/pkg/apis/ark/v1"
	"github.com/heptio/ark/pkg/backup"
	"github.com/heptio/ark/pkg/client"
	"github.com/heptio/ark/pkg/cloudprovider"
	arkaws "github.com/heptio/ark/pkg/cloudprovider/aws"
	"github.com/heptio/ark/pkg/cloudprovider/azure"
	"github.com/heptio/ark/pkg/cloudprovider/gcp"
	"github.com/heptio/ark/pkg/cmd"
	"github.com/heptio/ark/pkg/controller"
	arkdiscovery "github.com/heptio/ark/pkg/discovery"
	"github.com/heptio/ark/pkg/generated/clientset"
	arkv1client "github.com/heptio/ark/pkg/generated/clientset/typed/ark/v1"
	informers "github.com/heptio/ark/pkg/generated/informers/externalversions"
	"github.com/heptio/ark/pkg/restore"
	"github.com/heptio/ark/pkg/restore/restorers"
	"github.com/heptio/ark/pkg/util/kube"
)

func NewCommand() *cobra.Command {
	var kubeconfig string

	var command = &cobra.Command{
		Use:   "server",
		Short: "Run the ark server",
		Long:  "Run the ark server",
		Run: func(c *cobra.Command, args []string) {
			s, err := newServer(kubeconfig)
			cmd.CheckError(err)

			cmd.CheckError(s.run())
		},
	}

	command.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use to talk to the Kubernetes apiserver. If unset, try the environment variable KUBECONFIG, as well as in-cluster configuration")

	return command
}

type server struct {
	kubeClient            kubernetes.Interface
	arkClient             clientset.Interface
	backupService         cloudprovider.BackupService
	snapshotService       cloudprovider.SnapshotService
	discoveryClient       discovery.DiscoveryInterface
	clientPool            dynamic.ClientPool
	sharedInformerFactory informers.SharedInformerFactory
	ctx                   context.Context
	cancelFunc            context.CancelFunc
}

func newServer(kubeconfig string) (*server, error) {
	clientConfig, err := client.Config(kubeconfig)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	arkClient, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &server{
		kubeClient:            kubeClient,
		arkClient:             arkClient,
		discoveryClient:       arkClient.Discovery(),
		clientPool:            dynamic.NewDynamicClientPool(clientConfig),
		sharedInformerFactory: informers.NewSharedInformerFactory(arkClient, 0),
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	return s, nil
}

func (s *server) run() error {
	if err := s.ensureArkNamespace(); err != nil {
		return err
	}

	config, err := s.loadConfig()
	if err != nil {
		return err
	}
	applyConfigDefaults(config)

	s.watchConfig(config)

	if err := s.initBackupService(config); err != nil {
		return err
	}

	if err := s.initSnapshotService(config); err != nil {
		return err
	}

	if err := s.runControllers(config); err != nil {
		return err
	}

	return nil
}

func (s *server) ensureArkNamespace() error {
	glog.Infof("Ensuring %s namespace exists for backups", api.DefaultNamespace)
	defaultNamespace := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: api.DefaultNamespace,
		},
	}

	if created, err := kube.EnsureNamespaceExists(&defaultNamespace, s.kubeClient.CoreV1().Namespaces()); created {
		glog.Infof("Namespace created")
	} else if err != nil {
		return err
	}
	glog.Infof("Namespace already exists")
	return nil
}

func (s *server) loadConfig() (*api.Config, error) {
	glog.Infof("Retrieving Ark configuration")
	var (
		config *api.Config
		err    error
	)
	for {
		config, err = s.arkClient.ArkV1().Configs(api.DefaultNamespace).Get("default", metav1.GetOptions{})
		if err == nil {
			break
		}
		if !apierrors.IsNotFound(err) {
			glog.Errorf("error retrieving configuration: %v", err)
		}
		glog.Infof("Will attempt to retrieve configuration again in 5 seconds")
		time.Sleep(5 * time.Second)
	}
	glog.Infof("Successfully retrieved Ark configuration")
	return config, nil
}

const (
	defaultGCSyncPeriod       = 60 * time.Minute
	defaultBackupSyncPeriod   = 60 * time.Minute
	defaultScheduleSyncPeriod = time.Minute
)

var defaultResourcePriorities = []string{
	"namespaces",
	"persistentvolumes",
	"persistentvolumeclaims",
	"secrets",
	"configmaps",
}

func applyConfigDefaults(c *api.Config) {
	if c.GCSyncPeriod.Duration == 0 {
		c.GCSyncPeriod.Duration = defaultGCSyncPeriod
	}

	if c.BackupSyncPeriod.Duration == 0 {
		c.BackupSyncPeriod.Duration = defaultBackupSyncPeriod
	}

	if c.ScheduleSyncPeriod.Duration == 0 {
		c.ScheduleSyncPeriod.Duration = defaultScheduleSyncPeriod
	}

	if len(c.ResourcePriorities) == 0 {
		c.ResourcePriorities = defaultResourcePriorities
		glog.Infof("Using default resource priorities: %v", c.ResourcePriorities)
	} else {
		glog.Infof("Using resource priorities from config: %v", c.ResourcePriorities)
	}
}

// watchConfig adds an update event handler to the Config shared informer, invoking s.cancelFunc
// when it sees a change.
func (s *server) watchConfig(config *api.Config) {
	s.sharedInformerFactory.Ark().V1().Configs().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			updated := newObj.(*api.Config)

			if updated.Name != config.Name {
				glog.V(5).Infof("config watch channel received other config %q", updated.Name)
				return
			}

			if !reflect.DeepEqual(config, updated) {
				glog.Infof("Detected a config change. Gracefully shutting down")
				s.cancelFunc()
			}
		},
	})
}

func (s *server) initBackupService(config *api.Config) error {
	glog.Infof("Configuring cloud provider for backup service")
	cloud, err := initCloud(config.BackupStorageProvider.CloudProviderConfig, "backupStorageProvider")
	if err != nil {
		return err
	}
	s.backupService = cloudprovider.NewBackupService(cloud.ObjectStorage())
	return nil
}

func (s *server) initSnapshotService(config *api.Config) error {
	glog.Infof("Configuring cloud provider for snapshot service")
	cloud, err := initCloud(config.PersistentVolumeProvider, "persistentVolumeProvider")
	if err != nil {
		return err
	}
	s.snapshotService = cloudprovider.NewSnapshotService(cloud.BlockStorage())
	return nil
}

func initCloud(config api.CloudProviderConfig, field string) (cloudprovider.StorageAdapter, error) {
	var (
		cloud cloudprovider.StorageAdapter
		err   error
	)

	if config.AWS != nil {
		cloud, err = getAWSCloudProvider(config)
	}

	if config.GCP != nil {
		if cloud != nil {
			return nil, fmt.Errorf("you may only specify one of aws, gcp, or azure for %s", field)
		}
		cloud, err = getGCPCloudProvider(config)
	}

	if config.Azure != nil {
		if cloud != nil {
			return nil, fmt.Errorf("you may only specify one of aws, gcp, or azure for %s", field)
		}
		cloud, err = getAzureCloudProvider(config)
	}

	if err != nil {
		return nil, err
	}

	if cloud == nil {
		return nil, fmt.Errorf("you must specify one of aws, gcp, or azure for %s", field)
	}

	return cloud, err
}

func getAWSCloudProvider(cloudConfig api.CloudProviderConfig) (cloudprovider.StorageAdapter, error) {
	if cloudConfig.AWS == nil {
		return nil, errors.New("missing aws configuration in config file")
	}
	if cloudConfig.AWS.Region == "" {
		return nil, errors.New("missing region in aws configuration in config file")
	}
	if cloudConfig.AWS.AvailabilityZone == "" {
		return nil, errors.New("missing availabilityZone in aws configuration in config file")
	}

	awsConfig := aws.NewConfig().
		WithRegion(cloudConfig.AWS.Region).
		WithS3ForcePathStyle(cloudConfig.AWS.S3ForcePathStyle)

	if cloudConfig.AWS.S3Url != "" {
		awsConfig = awsConfig.WithEndpointResolver(
			endpoints.ResolverFunc(func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
				if service == endpoints.S3ServiceID {
					return endpoints.ResolvedEndpoint{
						URL: cloudConfig.AWS.S3Url,
					}, nil
				}

				return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
			}),
		)
	}

	return arkaws.NewStorageAdapter(awsConfig, cloudConfig.AWS.AvailabilityZone, cloudConfig.AWS.KMSKeyID)
}

func getGCPCloudProvider(cloudConfig api.CloudProviderConfig) (cloudprovider.StorageAdapter, error) {
	if cloudConfig.GCP == nil {
		return nil, errors.New("missing gcp configuration in config file")
	}
	if cloudConfig.GCP.Project == "" {
		return nil, errors.New("missing project in gcp configuration in config file")
	}
	if cloudConfig.GCP.Zone == "" {
		return nil, errors.New("missing zone in gcp configuration in config file")
	}
	return gcp.NewStorageAdapter(cloudConfig.GCP.Project, cloudConfig.GCP.Zone)
}

func getAzureCloudProvider(cloudConfig api.CloudProviderConfig) (cloudprovider.StorageAdapter, error) {
	if cloudConfig.Azure == nil {
		return nil, errors.New("missing azure configuration in config file")
	}
	if cloudConfig.Azure.Location == "" {
		return nil, errors.New("missing location in azure configuration in config file")
	}
	return azure.NewStorageAdapter(cloudConfig.Azure.Location, cloudConfig.Azure.APITimeout.Duration)
}

func durationMin(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func (s *server) runControllers(config *api.Config) error {
	glog.Infof("Starting controllers")

	ctx := s.ctx
	var wg sync.WaitGroup

	cloudBackupCacheResyncPeriod := durationMin(config.GCSyncPeriod.Duration, config.BackupSyncPeriod.Duration)
	glog.Infof("Caching cloud backups every %s", cloudBackupCacheResyncPeriod)
	s.backupService = cloudprovider.NewBackupServiceWithCachedBackupGetter(
		ctx,
		s.backupService,
		cloudBackupCacheResyncPeriod,
	)

	backupSyncController := controller.NewBackupSyncController(
		s.arkClient.ArkV1(),
		s.backupService,
		config.BackupStorageProvider.Bucket,
		config.BackupSyncPeriod.Duration,
	)
	wg.Add(1)
	go func() {
		backupSyncController.Run(ctx, 1)
		wg.Done()
	}()

	discoveryHelper, err := arkdiscovery.NewHelper(s.discoveryClient)
	if err != nil {
		return err
	}
	go wait.Until(
		func() {
			if err := discoveryHelper.Refresh(); err != nil {
				glog.Errorf("error refreshing discovery: %v", err)
			}
		},
		5*time.Minute,
		ctx.Done(),
	)

	if config.RestoreOnlyMode {
		glog.Infof("Restore only mode - not starting the backup, schedule or GC controllers")
	} else {
		backupper, err := newBackupper(discoveryHelper, s.clientPool, s.backupService, s.snapshotService)
		cmd.CheckError(err)
		backupController := controller.NewBackupController(
			s.sharedInformerFactory.Ark().V1().Backups(),
			s.arkClient.ArkV1(),
			backupper,
			s.backupService,
			config.BackupStorageProvider.Bucket,
			discoveryHelper.Mapper(),
		)
		wg.Add(1)
		go func() {
			backupController.Run(ctx, 1)
			wg.Done()
		}()

		scheduleController := controller.NewScheduleController(
			s.arkClient.ArkV1(),
			s.arkClient.ArkV1(),
			s.sharedInformerFactory.Ark().V1().Schedules(),
			config.ScheduleSyncPeriod.Duration,
		)
		wg.Add(1)
		go func() {
			scheduleController.Run(ctx, 1)
			wg.Done()
		}()

		gcController := controller.NewGCController(
			s.backupService,
			s.snapshotService,
			config.BackupStorageProvider.Bucket,
			config.GCSyncPeriod.Duration,
			s.sharedInformerFactory.Ark().V1().Backups(),
			s.arkClient.ArkV1(),
		)
		wg.Add(1)
		go func() {
			gcController.Run(ctx, 1)
			wg.Done()
		}()
	}

	restorer, err := newRestorer(
		discoveryHelper,
		s.clientPool,
		s.backupService,
		s.snapshotService,
		config.ResourcePriorities,
		s.arkClient.ArkV1(),
		s.kubeClient,
	)
	cmd.CheckError(err)

	restoreController := controller.NewRestoreController(
		s.sharedInformerFactory.Ark().V1().Restores(),
		s.arkClient.ArkV1(),
		s.arkClient.ArkV1(),
		restorer,
		s.backupService,
		config.BackupStorageProvider.Bucket,
		s.sharedInformerFactory.Ark().V1().Backups(),
	)
	wg.Add(1)
	go func() {
		restoreController.Run(ctx, 1)
		wg.Done()
	}()

	// SHARED INFORMERS HAVE TO BE STARTED AFTER ALL CONTROLLERS
	go s.sharedInformerFactory.Start(ctx.Done())

	glog.Infof("Server started successfully")

	<-ctx.Done()

	glog.Info("Waiting for all controllers to shut down gracefully")
	wg.Wait()

	return nil
}

func newBackupper(
	discoveryHelper arkdiscovery.Helper,
	clientPool dynamic.ClientPool,
	backupService cloudprovider.BackupService,
	snapshotService cloudprovider.SnapshotService,
) (backup.Backupper, error) {
	actions := map[string]backup.Action{}

	if snapshotService != nil {
		actions["persistentvolumes"] = backup.NewVolumeSnapshotAction(snapshotService)
	}

	return backup.NewKubernetesBackupper(
		discoveryHelper,
		client.NewDynamicFactory(clientPool),
		actions,
	)
}

func newRestorer(
	discoveryHelper arkdiscovery.Helper,
	clientPool dynamic.ClientPool,
	backupService cloudprovider.BackupService,
	snapshotService cloudprovider.SnapshotService,
	resourcePriorities []string,
	backupClient arkv1client.BackupsGetter,
	kubeClient kubernetes.Interface,
) (restore.Restorer, error) {
	restorers := map[string]restorers.ResourceRestorer{
		"persistentvolumes":      restorers.NewPersistentVolumeRestorer(snapshotService),
		"persistentvolumeclaims": restorers.NewPersistentVolumeClaimRestorer(),
		"services":               restorers.NewServiceRestorer(),
		"namespaces":             restorers.NewNamespaceRestorer(),
		"pods":                   restorers.NewPodRestorer(),
		"jobs":                   restorers.NewJobRestorer(),
	}

	return restore.NewKubernetesRestorer(
		discoveryHelper,
		client.NewDynamicFactory(clientPool),
		restorers,
		backupService,
		resourcePriorities,
		backupClient,
		kubeClient.CoreV1().Namespaces(),
	)
}
