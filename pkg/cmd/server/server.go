/*
Copyright 2017 the Heptio Ark contributors.

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	kubeerrs "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	api "github.com/heptio/ark/pkg/apis/ark/v1"
	"github.com/heptio/ark/pkg/backup"
	"github.com/heptio/ark/pkg/buildinfo"
	"github.com/heptio/ark/pkg/client"
	"github.com/heptio/ark/pkg/cloudprovider"
	"github.com/heptio/ark/pkg/cmd"
	"github.com/heptio/ark/pkg/cmd/util/signals"
	"github.com/heptio/ark/pkg/controller"
	arkdiscovery "github.com/heptio/ark/pkg/discovery"
	clientset "github.com/heptio/ark/pkg/generated/clientset/versioned"
	informers "github.com/heptio/ark/pkg/generated/informers/externalversions"
	"github.com/heptio/ark/pkg/metrics"
	"github.com/heptio/ark/pkg/plugin"
	"github.com/heptio/ark/pkg/podexec"
	"github.com/heptio/ark/pkg/restic"
	"github.com/heptio/ark/pkg/restore"
	"github.com/heptio/ark/pkg/util/kube"
	"github.com/heptio/ark/pkg/util/logging"
	"github.com/heptio/ark/pkg/util/stringslice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// the port where prometheus metrics are exposed
	defaultMetricsAddress = ":8085"
)

type serverConfig struct {
	pluginDir, metricsAddress, defaultBackupLocation string
	backupSyncPeriod, podVolumeOperationTimeout      time.Duration
	restoreResourcePriorities                        []string
	restoreOnly                                      bool
}

func NewCommand() *cobra.Command {
	var (
		logLevelFlag = logging.LogLevelFlag(logrus.InfoLevel)
		config       = serverConfig{
			pluginDir:                 "/plugins",
			metricsAddress:            defaultMetricsAddress,
			defaultBackupLocation:     "default",
			backupSyncPeriod:          defaultBackupSyncPeriod,
			podVolumeOperationTimeout: defaultPodVolumeOperationTimeout,
			restoreResourcePriorities: defaultRestorePriorities,
		}
	)

	var command = &cobra.Command{
		Use:   "server",
		Short: "Run the ark server",
		Long:  "Run the ark server",
		Run: func(c *cobra.Command, args []string) {
			// go-plugin uses log.Println to log when it's waiting for all plugin processes to complete so we need to
			// set its output to stdout.
			log.SetOutput(os.Stdout)

			logLevel := logLevelFlag.Parse()
			// Make sure we log to stdout so cloud log dashboards don't show this as an error.
			logrus.SetOutput(os.Stdout)
			logrus.Infof("setting log-level to %s", strings.ToUpper(logLevel.String()))

			// Ark's DefaultLogger logs to stdout, so all is good there.
			logger := logging.DefaultLogger(logLevel)
			logger.Infof("Starting Ark server %s", buildinfo.FormattedGitSHA())

			// NOTE: the namespace flag is bound to ark's persistent flags when the root ark command
			// creates the client Factory and binds the Factory's flags. We're not using a Factory here in
			// the server because the Factory gets its basename set at creation time, and the basename is
			// used to construct the user-agent for clients. Also, the Factory's Namespace() method uses
			// the client config file to determine the appropriate namespace to use, and that isn't
			// applicable to the server (it uses the method directly below instead). We could potentially
			// add a SetBasename() method to the Factory, and tweak how Namespace() works, if we wanted to
			// have the server use the Factory.
			namespaceFlag := c.Flag("namespace")
			if namespaceFlag == nil {
				cmd.CheckError(errors.New("unable to look up namespace flag"))
			}
			namespace := getServerNamespace(namespaceFlag)

			s, err := newServer(namespace, fmt.Sprintf("%s-%s", c.Parent().Name(), c.Name()), config, logger)
			cmd.CheckError(err)

			cmd.CheckError(s.run())
		},
	}

	command.Flags().Var(logLevelFlag, "log-level", fmt.Sprintf("the level at which to log. Valid values are %s.", strings.Join(logLevelFlag.AllowedValues(), ", ")))
	command.Flags().StringVar(&config.pluginDir, "plugin-dir", config.pluginDir, "directory containing Ark plugins")
	command.Flags().StringVar(&config.metricsAddress, "metrics-address", config.metricsAddress, "the address to expose prometheus metrics")
	command.Flags().DurationVar(&config.backupSyncPeriod, "backup-sync-period", config.backupSyncPeriod, "how often to ensure all Ark backups in object storage exist as Backup API objects in the cluster")
	command.Flags().DurationVar(&config.podVolumeOperationTimeout, "restic-timeout", config.podVolumeOperationTimeout, "how long backups/restores of pod volumes should be allowed to run before timing out")
	command.Flags().BoolVar(&config.restoreOnly, "restore-only", config.restoreOnly, "run in a mode where only restores are allowed; backups, schedules, and garbage-collection are all disabled")
	command.Flags().StringSliceVar(&config.restoreResourcePriorities, "restore-resource-priorities", config.restoreResourcePriorities, "desired order of resource restores; any resource not in the list will be restored alphabetically after the prioritized resources")
	command.Flags().StringVar(&config.defaultBackupLocation, "default-backup-storage-location", config.defaultBackupLocation, "name of the default backup storage location")

	return command
}

func getServerNamespace(namespaceFlag *pflag.Flag) string {
	if namespaceFlag.Changed {
		return namespaceFlag.Value.String()
	}

	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}

	return api.DefaultNamespace
}

type server struct {
	namespace             string
	metricsAddress        string
	kubeClientConfig      *rest.Config
	kubeClient            kubernetes.Interface
	arkClient             clientset.Interface
	blockStore            cloudprovider.BlockStore
	discoveryClient       discovery.DiscoveryInterface
	discoveryHelper       arkdiscovery.Helper
	dynamicClient         dynamic.Interface
	sharedInformerFactory informers.SharedInformerFactory
	ctx                   context.Context
	cancelFunc            context.CancelFunc
	logger                logrus.FieldLogger
	logLevel              logrus.Level
	pluginRegistry        plugin.Registry
	pluginManager         plugin.Manager
	resticManager         restic.RepositoryManager
	metrics               *metrics.ServerMetrics
	config                serverConfig
}

func newServer(namespace, baseName string, config serverConfig, logger *logrus.Logger) (*server, error) {
	clientConfig, err := client.Config("", "", baseName)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	arkClient, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pluginRegistry := plugin.NewRegistry(config.pluginDir, logger, logger.Level)
	if err := pluginRegistry.DiscoverPlugins(); err != nil {
		return nil, err
	}
	pluginManager := plugin.NewManager(logger, logger.Level, pluginRegistry)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &server{
		namespace:             namespace,
		metricsAddress:        config.metricsAddress,
		kubeClientConfig:      clientConfig,
		kubeClient:            kubeClient,
		arkClient:             arkClient,
		discoveryClient:       arkClient.Discovery(),
		dynamicClient:         dynamicClient,
		sharedInformerFactory: informers.NewFilteredSharedInformerFactory(arkClient, 0, namespace, nil),
		ctx:            ctx,
		cancelFunc:     cancelFunc,
		logger:         logger,
		logLevel:       logger.Level,
		pluginRegistry: pluginRegistry,
		pluginManager:  pluginManager,
		config:         config,
	}

	return s, nil
}

func (s *server) run() error {
	defer s.pluginManager.CleanupClients()

	signals.CancelOnShutdown(s.cancelFunc, s.logger)

	// Since s.namespace, which specifies where backups/restores/schedules/etc. should live,
	// *could* be different from the namespace where the Ark server pod runs, check to make
	// sure it exists, and fail fast if it doesn't.
	if err := s.namespaceExists(s.namespace); err != nil {
		return err
	}

	if err := s.initDiscoveryHelper(); err != nil {
		return err
	}

	// check to ensure all Ark CRDs exist
	if err := s.arkResourcesExist(); err != nil {
		return err
	}

	originalConfig, err := s.loadConfig()
	if err != nil {
		return err
	}

	// watchConfig needs to examine the unmodified original config, so we keep that around as a
	// separate object, and instead apply defaults to a clone.
	config := originalConfig.DeepCopy()
	s.applyConfigDefaults(config)

	s.watchConfig(originalConfig)

	backupStorageLocation, err := s.arkClient.ArkV1().BackupStorageLocations(s.namespace).Get(s.config.defaultBackupLocation, metav1.GetOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	if config.PersistentVolumeProvider == nil {
		s.logger.Info("PersistentVolumeProvider config not provided, volume snapshots and restores are disabled")
	} else {
		s.logger.Info("Configuring cloud provider for snapshot service")
		blockStore, err := getBlockStore(*config.PersistentVolumeProvider, s.pluginManager)
		if err != nil {
			return err
		}
		s.blockStore = blockStore
	}

	if backupStorageLocation.Spec.Config[restic.ResticLocationConfigKey] != "" {
		if err := s.initRestic(backupStorageLocation.Spec.Provider); err != nil {
			return err
		}
	}

	if err := s.runControllers(config, backupStorageLocation); err != nil {
		return err
	}

	return nil
}

func (s *server) applyConfigDefaults(c *api.Config) {
	if s.config.backupSyncPeriod == 0 {
		s.config.backupSyncPeriod = defaultBackupSyncPeriod
	}

	if s.config.podVolumeOperationTimeout == 0 {
		s.config.podVolumeOperationTimeout = defaultPodVolumeOperationTimeout
	}

	if len(s.config.restoreResourcePriorities) == 0 {
		s.config.restoreResourcePriorities = defaultRestorePriorities
		s.logger.WithField("priorities", s.config.restoreResourcePriorities).Info("Using default resource priorities")
	} else {
		s.logger.WithField("priorities", s.config.restoreResourcePriorities).Info("Using given resource priorities")
	}
}

// namespaceExists returns nil if namespace can be successfully
// gotten from the kubernetes API, or an error otherwise.
func (s *server) namespaceExists(namespace string) error {
	s.logger.WithField("namespace", namespace).Info("Checking existence of namespace")

	if _, err := s.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{}); err != nil {
		return errors.WithStack(err)
	}

	s.logger.WithField("namespace", namespace).Info("Namespace exists")
	return nil
}

// initDiscoveryHelper instantiates the server's discovery helper and spawns a
// goroutine to call Refresh() every 5 minutes.
func (s *server) initDiscoveryHelper() error {
	discoveryHelper, err := arkdiscovery.NewHelper(s.discoveryClient, s.logger)
	if err != nil {
		return err
	}
	s.discoveryHelper = discoveryHelper

	go wait.Until(
		func() {
			if err := discoveryHelper.Refresh(); err != nil {
				s.logger.WithError(err).Error("Error refreshing discovery")
			}
		},
		5*time.Minute,
		s.ctx.Done(),
	)

	return nil
}

// arkResourcesExist checks for the existence of each Ark CRD via discovery
// and returns an error if any of them don't exist.
func (s *server) arkResourcesExist() error {
	s.logger.Info("Checking existence of Ark custom resource definitions")

	var arkGroupVersion *metav1.APIResourceList
	for _, gv := range s.discoveryHelper.Resources() {
		if gv.GroupVersion == api.SchemeGroupVersion.String() {
			arkGroupVersion = gv
			break
		}
	}

	if arkGroupVersion == nil {
		return errors.Errorf("Ark API group %s not found. Apply examples/common/00-prereqs.yaml to create it.", api.SchemeGroupVersion)
	}

	foundResources := sets.NewString()
	for _, resource := range arkGroupVersion.APIResources {
		foundResources.Insert(resource.Kind)
	}

	var errs []error
	for kind := range api.CustomResources() {
		if foundResources.Has(kind) {
			s.logger.WithField("kind", kind).Debug("Found custom resource")
			continue
		}

		errs = append(errs, errors.Errorf("custom resource %s not found in Ark API group %s", kind, api.SchemeGroupVersion))
	}

	if len(errs) > 0 {
		errs = append(errs, errors.New("Ark custom resources not found - apply examples/common/00-prereqs.yaml to update the custom resource definitions"))
		return kubeerrs.NewAggregate(errs)
	}

	s.logger.Info("All Ark custom resource definitions exist")
	return nil
}

func (s *server) loadConfig() (*api.Config, error) {
	s.logger.Info("Retrieving Ark configuration")
	var (
		config *api.Config
		err    error
	)
	for {
		config, err = s.arkClient.ArkV1().Configs(s.namespace).Get("default", metav1.GetOptions{})
		if err == nil {
			break
		}
		if !apierrors.IsNotFound(err) {
			s.logger.WithError(err).Error("Error retrieving configuration")
		} else {
			s.logger.Info("Configuration not found")
		}
		s.logger.Info("Will attempt to retrieve configuration again in 5 seconds")
		time.Sleep(5 * time.Second)
	}
	s.logger.Info("Successfully retrieved Ark configuration")
	return config, nil
}

const (
	defaultBackupSyncPeriod          = 60 * time.Minute
	defaultPodVolumeOperationTimeout = 60 * time.Minute
)

// - Namespaces go first because all namespaced resources depend on them.
// - PVs go before PVCs because PVCs depend on them.
// - PVCs go before pods or controllers so they can be mounted as volumes.
// - Secrets and config maps go before pods or controllers so they can be mounted
// 	 as volumes.
// - Service accounts go before pods or controllers so pods can use them.
// - Limit ranges go before pods or controllers so pods can use them.
// - Pods go before controllers so they can be explicitly restored and potentially
//	 have restic restores run before controllers adopt the pods.
var defaultRestorePriorities = []string{
	"namespaces",
	"persistentvolumes",
	"persistentvolumeclaims",
	"secrets",
	"configmaps",
	"serviceaccounts",
	"limitranges",
	"pods",
}

// watchConfig adds an update event handler to the Config shared informer, invoking s.cancelFunc
// when it sees a change.
func (s *server) watchConfig(config *api.Config) {
	s.sharedInformerFactory.Ark().V1().Configs().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			updated := newObj.(*api.Config)
			s.logger.WithField("name", kube.NamespaceAndName(updated)).Debug("received updated config")

			if updated.Name != config.Name {
				s.logger.WithField("name", updated.Name).Debug("Config watch channel received other config")
				return
			}

			// Objects retrieved via Get() don't have their Kind or APIVersion set. Objects retrieved via
			// Watch(), including those from shared informer event handlers, DO have their Kind and
			// APIVersion set. To prevent the DeepEqual() call below from considering Kind or APIVersion
			// as the source of a change, set config.Kind and config.APIVersion to match the values from
			// the updated Config.
			if config.Kind != updated.Kind {
				config.Kind = updated.Kind
			}
			if config.APIVersion != updated.APIVersion {
				config.APIVersion = updated.APIVersion
			}

			if !reflect.DeepEqual(config, updated) {
				s.logger.Info("Detected a config change. Gracefully shutting down")
				s.cancelFunc()
			}
		},
	})
}

func getBlockStore(cloudConfig api.CloudProviderConfig, manager plugin.Manager) (cloudprovider.BlockStore, error) {
	if cloudConfig.Name == "" {
		return nil, errors.New("block storage provider name must not be empty")
	}

	blockStore, err := manager.GetBlockStore(cloudConfig.Name)
	if err != nil {
		return nil, err
	}

	if err := blockStore.Init(cloudConfig.Config); err != nil {
		return nil, err
	}

	return blockStore, nil
}

func durationMin(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func (s *server) initRestic(providerName string) error {
	// warn if restic daemonset does not exist
	if _, err := s.kubeClient.AppsV1().DaemonSets(s.namespace).Get(restic.DaemonSet, metav1.GetOptions{}); apierrors.IsNotFound(err) {
		s.logger.Warn("Ark restic daemonset not found; restic backups/restores will not work until it's created")
	} else if err != nil {
		s.logger.WithError(errors.WithStack(err)).Warn("Error checking for existence of ark restic daemonset")
	}

	// ensure the repo key secret is set up
	if err := restic.EnsureCommonRepositoryKey(s.kubeClient.CoreV1(), s.namespace); err != nil {
		return err
	}

	// set the env vars that restic uses for creds purposes
	if providerName == string(restic.AzureBackend) {
		os.Setenv("AZURE_ACCOUNT_NAME", os.Getenv("AZURE_STORAGE_ACCOUNT_ID"))
		os.Setenv("AZURE_ACCOUNT_KEY", os.Getenv("AZURE_STORAGE_KEY"))
	}

	// use a stand-alone secrets informer so we can filter to only the restic credentials
	// secret(s) within the heptio-ark namespace
	//
	// note: using an informer to access the single secret for all ark-managed
	// restic repositories is overkill for now, but will be useful when we move
	// to fully-encrypted backups and have unique keys per repository.
	secretsInformer := corev1informers.NewFilteredSecretInformer(
		s.kubeClient,
		s.namespace,
		0,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
		func(opts *metav1.ListOptions) {
			opts.FieldSelector = fmt.Sprintf("metadata.name=%s", restic.CredentialsSecretName)
		},
	)
	go secretsInformer.Run(s.ctx.Done())

	res, err := restic.NewRepositoryManager(
		s.ctx,
		s.namespace,
		s.arkClient,
		secretsInformer,
		s.sharedInformerFactory.Ark().V1().ResticRepositories(),
		s.arkClient.ArkV1(),
		s.logger,
	)
	if err != nil {
		return err
	}
	s.resticManager = res

	return nil
}

func (s *server) runControllers(config *api.Config, defaultBackupLocation *api.BackupStorageLocation) error {
	s.logger.Info("Starting controllers")

	ctx := s.ctx
	var wg sync.WaitGroup

	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		s.logger.Infof("Starting metric server at address [%s]", s.metricsAddress)
		if err := http.ListenAndServe(s.metricsAddress, metricsMux); err != nil {
			s.logger.Fatalf("Failed to start metric server at [%s]: %v", s.metricsAddress, err)
		}
	}()
	s.metrics = metrics.NewServerMetrics()
	s.metrics.RegisterAllMetrics()

	backupSyncController := controller.NewBackupSyncController(
		s.arkClient.ArkV1(),
		s.sharedInformerFactory.Ark().V1().Backups(),
		s.sharedInformerFactory.Ark().V1().BackupStorageLocations(),
		s.config.backupSyncPeriod,
		s.namespace,
		s.pluginRegistry,
		s.logger,
		s.logLevel,
	)
	wg.Add(1)
	go func() {
		backupSyncController.Run(ctx, 1)
		wg.Done()
	}()

	if s.config.restoreOnly {
		s.logger.Info("Restore only mode - not starting the backup, schedule, delete-backup, or GC controllers")
	} else {
		backupTracker := controller.NewBackupTracker()

		backupper, err := backup.NewKubernetesBackupper(
			s.discoveryHelper,
			client.NewDynamicFactory(s.dynamicClient),
			podexec.NewPodCommandExecutor(s.kubeClientConfig, s.kubeClient.CoreV1().RESTClient()),
			s.blockStore,
			s.resticManager,
			s.config.podVolumeOperationTimeout,
		)
		cmd.CheckError(err)

		backupController := controller.NewBackupController(
			s.sharedInformerFactory.Ark().V1().Backups(),
			s.arkClient.ArkV1(),
			backupper,
			s.blockStore != nil,
			s.logger,
			s.logLevel,
			s.pluginRegistry,
			backupTracker,
			s.sharedInformerFactory.Ark().V1().BackupStorageLocations(),
			s.config.defaultBackupLocation,
			s.metrics,
		)
		wg.Add(1)
		go func() {
			backupController.Run(ctx, 1)
			wg.Done()
		}()

		scheduleController := controller.NewScheduleController(
			s.namespace,
			s.arkClient.ArkV1(),
			s.arkClient.ArkV1(),
			s.sharedInformerFactory.Ark().V1().Schedules(),
			s.logger,
			s.metrics,
		)
		wg.Add(1)
		go func() {
			scheduleController.Run(ctx, 1)
			wg.Done()
		}()

		gcController := controller.NewGCController(
			s.logger,
			s.sharedInformerFactory.Ark().V1().Backups(),
			s.sharedInformerFactory.Ark().V1().DeleteBackupRequests(),
			s.arkClient.ArkV1(),
		)
		wg.Add(1)
		go func() {
			gcController.Run(ctx, 1)
			wg.Done()
		}()

		backupDeletionController := controller.NewBackupDeletionController(
			s.logger,
			s.logLevel,
			s.sharedInformerFactory.Ark().V1().DeleteBackupRequests(),
			s.arkClient.ArkV1(), // deleteBackupRequestClient
			s.arkClient.ArkV1(), // backupClient
			s.blockStore,
			s.sharedInformerFactory.Ark().V1().Restores(),
			s.arkClient.ArkV1(), // restoreClient
			backupTracker,
			s.resticManager,
			s.sharedInformerFactory.Ark().V1().PodVolumeBackups(),
			s.sharedInformerFactory.Ark().V1().BackupStorageLocations(),
			s.pluginRegistry,
		)
		wg.Add(1)
		go func() {
			backupDeletionController.Run(ctx, 1)
			wg.Done()
		}()

	}

	restorer, err := restore.NewKubernetesRestorer(
		s.discoveryHelper,
		client.NewDynamicFactory(s.dynamicClient),
		s.blockStore,
		s.config.restoreResourcePriorities,
		s.arkClient.ArkV1(),
		s.kubeClient.CoreV1().Namespaces(),
		s.resticManager,
		s.config.podVolumeOperationTimeout,
		s.logger,
	)
	cmd.CheckError(err)

	restoreController := controller.NewRestoreController(
		s.namespace,
		s.sharedInformerFactory.Ark().V1().Restores(),
		s.arkClient.ArkV1(),
		s.arkClient.ArkV1(),
		restorer,
		s.sharedInformerFactory.Ark().V1().Backups(),
		s.sharedInformerFactory.Ark().V1().BackupStorageLocations(),
		s.blockStore != nil,
		s.logger,
		s.logLevel,
		s.pluginRegistry,
		s.config.defaultBackupLocation,
		s.metrics,
	)

	wg.Add(1)
	go func() {
		restoreController.Run(ctx, 1)
		wg.Done()
	}()

	downloadRequestController := controller.NewDownloadRequestController(
		s.arkClient.ArkV1(),
		s.sharedInformerFactory.Ark().V1().DownloadRequests(),
		s.sharedInformerFactory.Ark().V1().Restores(),
		s.sharedInformerFactory.Ark().V1().BackupStorageLocations(),
		s.sharedInformerFactory.Ark().V1().Backups(),
		s.pluginRegistry,
		s.logger,
		s.logLevel,
	)
	wg.Add(1)
	go func() {
		downloadRequestController.Run(ctx, 1)
		wg.Done()
	}()

	if s.resticManager != nil {
		resticRepoController := controller.NewResticRepositoryController(
			s.logger,
			s.sharedInformerFactory.Ark().V1().ResticRepositories(),
			s.arkClient.ArkV1(),
			defaultBackupLocation,
			s.resticManager,
		)
		wg.Add(1)
		go func() {
			// TODO only having a single worker may be an issue since maintenance
			// can take a long time.
			resticRepoController.Run(ctx, 1)
			wg.Done()
		}()
	}

	// SHARED INFORMERS HAVE TO BE STARTED AFTER ALL CONTROLLERS
	go s.sharedInformerFactory.Start(ctx.Done())

	// TODO(1.0): remove
	cache.WaitForCacheSync(ctx.Done(), s.sharedInformerFactory.Ark().V1().Backups().Informer().HasSynced)
	s.removeDeprecatedGCFinalizer()

	s.logger.Info("Server started successfully")

	<-ctx.Done()

	s.logger.Info("Waiting for all controllers to shut down gracefully")
	wg.Wait()

	return nil
}

// TODO(1.0): remove
func (s *server) removeDeprecatedGCFinalizer() {
	const gcFinalizer = "gc.ark.heptio.com"

	backups, err := s.sharedInformerFactory.Ark().V1().Backups().Lister().List(labels.Everything())
	if err != nil {
		s.logger.WithError(errors.WithStack(err)).Error("error listing backups from cache - unable to remove old finalizers")
		return
	}

	for _, backup := range backups {
		log := s.logger.WithField("backup", kube.NamespaceAndName(backup))

		if !stringslice.Has(backup.Finalizers, gcFinalizer) {
			log.Debug("backup doesn't have deprecated finalizer - skipping")
			continue
		}

		log.Info("removing deprecated finalizer from backup")

		patch := map[string]interface{}{
			"metadata": map[string]interface{}{
				"finalizers":      stringslice.Except(backup.Finalizers, gcFinalizer),
				"resourceVersion": backup.ResourceVersion,
			},
		}

		patchBytes, err := json.Marshal(patch)
		if err != nil {
			log.WithError(errors.WithStack(err)).Error("error marshaling finalizers patch")
			continue
		}

		_, err = s.arkClient.ArkV1().Backups(backup.Namespace).Patch(backup.Name, types.MergePatchType, patchBytes)
		if err != nil {
			log.WithError(errors.WithStack(err)).Error("error marshaling finalizers patch")
		}
	}
}
