/*
Copyright The Velero contributors.

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

package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	corev1listers "k8s.io/client-go/listers/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	listers "github.com/vmware-tanzu/velero/pkg/generated/listers/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/metrics"
	"github.com/vmware-tanzu/velero/pkg/restic"
	"github.com/vmware-tanzu/velero/pkg/util/filesystem"
	"github.com/vmware-tanzu/velero/pkg/util/kube"
)

// PodVolumeBackupReconciler reconciles a PodVolumeBackup object
type PodVolumeBackupReconciler struct {
	Scheme     *runtime.Scheme
	Client     client.Client
	Ctx        context.Context
	Clock      clock.Clock
	Metrics    *metrics.ServerMetrics
	NodeName   string
	FileSystem filesystem.Interface
	Log        logrus.FieldLogger

	PvLister  corev1listers.PersistentVolumeLister
	PvcLister corev1listers.PersistentVolumeClaimLister
	PvbLister listers.PodVolumeBackupLister
}

// +kubebuilder:rbac:groups=velero.io,resources=podvolumebackup,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=velero.io,resources=podvolumebackup/status,verbs=get;update;patch
func (r *PodVolumeBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithFields(logrus.Fields{
		"controller":      "podvolumebackup",
		"podvolumebackup": req.NamespacedName,
	})

	pvb := velerov1api.PodVolumeBackup{}
	if err := r.Client.Get(r.Ctx, req.NamespacedName, &pvb); err != nil {
		if apierrors.IsNotFound(err) {
			log.Debug("Unable to find PodVolumeBackup")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, errors.Wrap(err, "getting PodVolumeBackup")
	}

	if len(pvb.OwnerReferences) == 1 {
		log = log.WithField(
			"backup",
			fmt.Sprintf("%s/%s", req.Namespace, pvb.OwnerReferences[0].Name),
		)
	}

	log.Info("Pod volume backup starting")

	// only process items for this node
	if pvb.Spec.Node != r.NodeName {
		return ctrl.Result{}, nil
	}

	// only process new items
	if pvb.Status.Phase != "" && pvb.Status.Phase != velerov1api.PodVolumeBackupPhaseNew {
		log.Debug("Pod volume backup is not new, not processing")
		return ctrl.Result{}, nil
	}

	r.Metrics.RegisterPodVolumeBackupEnqueue(r.NodeName)

	// update status to InProgress
	pvb.Status.Phase = velerov1api.PodVolumeBackupPhaseInProgress
	pvb.Status.StartTimestamp = &metav1.Time{Time: r.Clock.Now()}

	pod := corev1.Pod{}
	if err := r.Client.Get(r.Ctx, req.NamespacedName, &pod); err != nil {
		return r.logErrorAndUpdateStatus(
			log,
			&pvb,
			err,
			fmt.Sprintf("getting pod %s/%s", pvb.Spec.Pod.Namespace, pvb.Spec.Pod.Name),
		)
	}

	volDir, err := kube.GetVolumeDirectory(&pod, pvb.Spec.Volume, r.PvcLister, r.PvLister)
	if err != nil {
		return r.logErrorAndUpdateStatus(log, &pvb, err, "getting volume directory name")
	}

	pathGlob := fmt.Sprintf("/host_pods/%s/volumes/*/%s", string(pvb.Spec.Pod.UID), volDir)
	log.WithField("pathGlob", pathGlob).Debug("Looking for path matching glob")

	path, err := singlePathMatch(pathGlob)
	if err != nil {
		return r.logErrorAndUpdateStatus(log, &pvb, err, "identifying unique volume path on host")
	}
	log.WithField("path", path).Debugf("Found path matching glob")

	// temporary credentials
	credsFile, err := restic.TempCredentialsFile(r.Client, req.Namespace, r.FileSystem)
	if err != nil {
		return r.logErrorAndUpdateStatus(log, &pvb, err, "creating temporary restic credentials file")
	}
	defer os.Remove(credsFile)

	resticCmd := restic.BackupCommand(pvb.Spec.RepoIdentifier, credsFile, path, pvb.Spec.Tags)

	caCert, err := restic.GetCACert(r.Client, pvb.Namespace, pvb.Spec.BackupStorageLocation)
	if err != nil {
		log.WithError(err).Error("getting caCert")
	}

	var caCertFile string
	if caCert != nil {
		caCertFile, err = restic.TempCACertFile(caCert, pvb.Spec.BackupStorageLocation, r.FileSystem)
		if err != nil {
			log.WithError(err).Error("creating temporary caCert file")
		}
		defer os.Remove(caCertFile)

	}

	resticCmd.CACertFile = caCertFile

	// Running restic command might need additional provider specific environment
	// variables. Based on the provider, we set resticCmd.Env appropriately (currently
	// for Azure and S3 based backuplocations)
	var env []string
	switch {
	case strings.HasPrefix(pvb.Spec.RepoIdentifier, "azure"):
		env, err = restic.AzureCmdEnv(r.Client, req.Namespace, pvb.Spec.BackupStorageLocation)
		if err != nil {
			return r.logErrorAndUpdateStatus(log, &pvb, err, "error setting restic with azure cmd env")

		}
	case strings.HasPrefix(pvb.Spec.RepoIdentifier, "s3"):
		env, err = restic.S3CmdEnv(r.Client, req.Namespace, pvb.Spec.BackupStorageLocation)
		if err != nil {
			return r.logErrorAndUpdateStatus(log, &pvb, err, "error setting restic with s3 cmd env")

		}
	}
	resticCmd.Env = env

	// If this is a PVC, look for the most recent completed pod volume backup for it and get
	// its restic snapshot ID to use as the value of the `--parent` flag. Without this,
	// if the pod using the PVC (and therefore the directory path under /host_pods/) has
	// changed since the PVC's last backup, restic will not be able to identify a suitable
	// parent snapshot to use, and will have to do a full rescan of the contents of the PVC.
	if pvcUID, ok := pvb.Labels[velerov1api.PVCUIDLabel]; ok {
		parentSnapshotID := getParentSnapshot(
			log,
			pvcUID,
			pvb.Spec.BackupStorageLocation,
			r.PvbLister.PodVolumeBackups(pvb.Namespace),
		)

		if parentSnapshotID == "" {
			log.Info("No parent snapshot found for PVC, not using --parent flag for this backup")
		} else {
			log.WithField("parentSnapshotID", parentSnapshotID).Info("Setting --parent flag for this backup")
			resticCmd.ExtraFlags = append(resticCmd.ExtraFlags, fmt.Sprintf("--parent=%s", parentSnapshotID))
		}
	}

	var emptySnapshot bool
	stdout, stderr, err := restic.RunBackup(resticCmd, log, r.updateBackupProgressFunc(&pvb, log))
	if err != nil {
		if strings.Contains(stderr, "snapshot is empty") {
			emptySnapshot = true
		} else {
			log.WithError(errors.WithStack(err)).Errorf(
				"Error running command=%s, stdout=%s, stderr=%s",
				resticCmd.String(),
				stdout,
				stderr,
			)
			return r.logErrorAndUpdateStatus(
				log,
				&pvb,
				err,
				fmt.Sprintf("running restic backup, stderr=%s", stderr),
			)
		}
	}
	log.Debugf("Ran command=%s, stdout=%s, stderr=%s", resticCmd.String(), stdout, stderr)

	var snapshotID string
	if !emptySnapshot {
		snapshotID, err = restic.GetSnapshotID(
			pvb.Spec.RepoIdentifier,
			credsFile,
			pvb.Spec.Tags,
			env,
			caCertFile,
		)
		if err != nil {
			return r.logErrorAndUpdateStatus(log, &pvb, err, "getting snapshot id")

		}
	}

	// update status to Completed with path & snapshot id
	pvb.Status.Path = path
	pvb.Status.Phase = velerov1api.PodVolumeBackupPhaseCompleted
	pvb.Status.SnapshotID = snapshotID
	pvb.Status.CompletionTimestamp = &metav1.Time{Time: r.Clock.Now()}
	if emptySnapshot {
		pvb.Status.Message = "volume was empty so no snapshot was taken"
	}

	latencyDuration := pvb.Status.CompletionTimestamp.Time.Sub(pvb.Status.StartTimestamp.Time)
	latencySeconds := float64(latencyDuration / time.Second)
	backupName := fmt.Sprintf("%s/%s", req.Namespace, pvb.OwnerReferences[0].Name)
	r.Metrics.ObserveResticOpLatency(r.NodeName, req.Name, resticCmd.Command, backupName, latencySeconds)
	r.Metrics.RegisterResticOpLatencyGauge(r.NodeName, req.Name, resticCmd.Command, backupName, latencySeconds)
	r.Metrics.RegisterPodVolumeBackupDequeue(r.NodeName)
	log.Info("Pod volume backup completed")

	return ctrl.Result{}, nil
}

func (r *PodVolumeBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&velerov1api.PodVolumeBackup{}).
		Complete(r)
}

func singlePathMatch(path string) (string, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if len(matches) != 1 {
		return "", errors.Errorf("expected one matching path, got %d", len(matches))
	}

	return matches[0], nil
}

// getParentSnapshot finds the most recent completed pod volume backup for the
// specified PVC and returns its restic snapshot ID. Any errors encountered are
// logged but not returned since they do not prevent a backup from proceeding.
func getParentSnapshot(log logrus.FieldLogger, pvcUID, backupStorageLocation string, pvbListers listers.PodVolumeBackupNamespaceLister) string {
	log = log.WithField("pvcUID", pvcUID)
	log.Infof("Looking for most recent completed pod volume backup for this PVC")

	pvcBackups, err := pvbListers.List(labels.SelectorFromSet(map[string]string{velerov1api.PVCUIDLabel: pvcUID}))
	if err != nil {
		log.WithError(errors.WithStack(err)).Error("Error listing pod volume backups for PVC")
		return ""
	}

	// go through all the pod volume backups for the PVC and look for the most recent completed one
	// to use as the parent.
	var mostRecentBackup *velerov1api.PodVolumeBackup
	for _, backup := range pvcBackups {
		if backup.Status.Phase != velerov1api.PodVolumeBackupPhaseCompleted {
			continue
		}

		if backupStorageLocation != backup.Spec.BackupStorageLocation {
			// Check the backup storage location is the same as spec in order to
			// support backup to multiple backup-locations. Otherwise, there exists
			// a case that backup volume snapshot to the second location would
			// failed, since the founded parent ID is only valid for the first
			// backup location, not the second backup location. Also, the second
			// backup should not use the first backup parent ID since its for the
			// first backup location only.
			continue
		}

		if mostRecentBackup == nil || backup.Status.StartTimestamp.After(mostRecentBackup.Status.StartTimestamp.Time) {
			mostRecentBackup = backup
		}
	}

	if mostRecentBackup == nil {
		log.Info("No completed pod volume backup found for PVC")
		return ""
	}

	log.WithFields(map[string]interface{}{
		"parentPodVolumeBackup": mostRecentBackup.Name,
		"parentSnapshotID":      mostRecentBackup.Status.SnapshotID,
	}).Info("Found most recent completed pod volume backup for PVC")

	return mostRecentBackup.Status.SnapshotID
}

// updateBackupProgressFunc returns a func that takes progress info and patches
// the PVB with the new progress
func (r *PodVolumeBackupReconciler) updateBackupProgressFunc(pvb *velerov1api.PodVolumeBackup, log logrus.FieldLogger) func(velerov1api.PodVolumeOperationProgress) {
	return func(progress velerov1api.PodVolumeOperationProgress) {
		pvb.Status.Progress = progress
	}
}

func (r *PodVolumeBackupReconciler) logErrorAndUpdateStatus(log *logrus.Entry, pvb *velerov1api.PodVolumeBackup, err error, msg string) (ctrl.Result, error) {
	log.WithError(err).Error(msg)
	pvb.Status.Phase = velerov1api.PodVolumeBackupPhaseFailed
	pvb.Status.Message = ""
	pvb.Status.CompletionTimestamp = &metav1.Time{Time: r.Clock.Now()}
	return ctrl.Result{}, errors.Wrap(err, msg)
}
