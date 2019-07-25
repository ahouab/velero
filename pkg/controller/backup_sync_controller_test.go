/*
Copyright 2017 the Velero contributors.

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
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	core "k8s.io/client-go/testing"

	velerov1api "github.com/heptio/velero/pkg/apis/velero/v1"
	pkgbackup "github.com/heptio/velero/pkg/backup"
	"github.com/heptio/velero/pkg/generated/clientset/versioned/fake"
	informers "github.com/heptio/velero/pkg/generated/informers/externalversions"
	"github.com/heptio/velero/pkg/label"
	"github.com/heptio/velero/pkg/persistence"
	persistencemocks "github.com/heptio/velero/pkg/persistence/mocks"
	"github.com/heptio/velero/pkg/plugin/clientmgmt"
	pluginmocks "github.com/heptio/velero/pkg/plugin/mocks"
	velerotest "github.com/heptio/velero/pkg/util/test"
)

func defaultPodVolumeBackup() *pkgbackup.PodVolumeBackupBuilder {
	return pkgbackup.NewNamedPodVolumeBackupBuilder(velerov1api.DefaultNamespace, "pvb-1")
}

func defaultLocationsList(namespace string) []*velerov1api.BackupStorageLocation {
	return []*velerov1api.BackupStorageLocation{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "location-1",
			},
			Spec: velerov1api.BackupStorageLocationSpec{
				Provider: "objStoreProvider",
				StorageType: velerov1api.StorageType{
					ObjectStorage: &velerov1api.ObjectStorageLocation{
						Bucket: "bucket-1",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "location-2",
			},
			Spec: velerov1api.BackupStorageLocationSpec{
				Provider: "objStoreProvider",
				StorageType: velerov1api.StorageType{
					ObjectStorage: &velerov1api.ObjectStorageLocation{
						Bucket: "bucket-2",
					},
				},
			},
		},
	}
}

func defaultLocationsListWithLongerLocationName(namespace string) []*velerov1api.BackupStorageLocation {
	return []*velerov1api.BackupStorageLocation{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "the-really-long-location-name-that-is-much-more-than-63-characters-1",
			},
			Spec: velerov1api.BackupStorageLocationSpec{
				Provider: "objStoreProvider",
				StorageType: velerov1api.StorageType{
					ObjectStorage: &velerov1api.ObjectStorageLocation{
						Bucket: "bucket-1",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "the-really-long-location-name-that-is-much-more-than-63-characters-2",
			},
			Spec: velerov1api.BackupStorageLocationSpec{
				Provider: "objStoreProvider",
				StorageType: velerov1api.StorageType{
					ObjectStorage: &velerov1api.ObjectStorageLocation{
						Bucket: "bucket-2",
					},
				},
			},
		},
	}
}

func TestBackupSyncControllerRun(t *testing.T) {
	type cloudBackupData struct {
		backup           *velerov1api.Backup
		podVolumeBackups []*velerov1api.PodVolumeBackup
	}

	tests := []struct {
		name                     string
		namespace                string
		locations                []*velerov1api.BackupStorageLocation
		cloudBuckets             map[string][]*cloudBackupData
		existingBackups          []*velerov1api.Backup
		existingPodVolumeBackups []*velerov1api.PodVolumeBackup
		longLocationNameEnabled  bool
	}{
		{
			name: "no cloud backups",
		},
		{
			name:      "normal case",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").Backup(),
					},
				},
			},
		},
		{
			name:      "all synced backups get created in Velero server's namespace",
			namespace: "velero",
			locations: defaultLocationsList("velero"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-2").Name("backup-3").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("velero").Name("backup-4").Backup(),
					},
				},
			},
		},
		{
			name:      "new backups get synced when some cloud backups already exist in the cluster",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-4").Backup(),
					},
				},
			},
			existingBackups: []*velerov1api.Backup{
				// add a label to each existing backup so we can differentiate it from the cloud
				// backup during verification
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels("i-exist", "true").StorageLocation("location-1").Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-3").Labels("i-exist", "true").StorageLocation("location-2").Backup(),
			},
		},
		{
			name:      "existing backups without a StorageLocation get it filled in",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
					},
				},
			},
			existingBackups: []*velerov1api.Backup{
				// add a label to each existing backup so we can differentiate it from the cloud
				// backup during verification
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels("i-exist", "true").StorageLocation("location-1").Backup(),
			},
		},
		{
			name:      "backup storage location names and labels get updated",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").StorageLocation("foo").Labels(velerov1api.StorageLocationLabel, "foo").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").StorageLocation("bar").Labels(velerov1api.StorageLocationLabel, "bar").Backup(),
					},
				},
			},
		},
		{
			name:                    "backup storage location names and labels get updated with location name greater than 63 chars",
			namespace:               "ns-1",
			locations:               defaultLocationsListWithLongerLocationName("ns-1"),
			longLocationNameEnabled: true,
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").StorageLocation("foo").Labels(velerov1api.StorageLocationLabel, "foo").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").StorageLocation("bar").Labels(velerov1api.StorageLocationLabel, "bar").Backup(),
					},
				},
			},
		},
		{
			name:      "all synced backups and pod volume backups get created in Velero server's namespace",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-1").PodVolumeBackup(),
						},
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-2").PodVolumeBackup(),
						},
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-4").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-1").PodVolumeBackup(),
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-2").PodVolumeBackup(),
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-3").PodVolumeBackup(),
						},
					},
				},
			},
		},
		{
			name:      "new pod volume backups get synched when some pod volume backups already exist in the cluster",
			namespace: "ns-1",
			locations: defaultLocationsList("ns-1"),
			cloudBuckets: map[string][]*cloudBackupData{
				"bucket-1": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-1").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-1").PodVolumeBackup(),
						},
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-2").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-3").PodVolumeBackup(),
						},
					},
				},
				"bucket-2": {
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-3").Backup(),
					},
					&cloudBackupData{
						backup: defaultBackup().Namespace("ns-1").Name("backup-4").Backup(),
						podVolumeBackups: []*velerov1api.PodVolumeBackup{
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-1").PodVolumeBackup(),
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-5").PodVolumeBackup(),
							defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-6").PodVolumeBackup(),
						},
					},
				},
			},
			existingPodVolumeBackups: []*velerov1api.PodVolumeBackup{
				defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-1").PodVolumeBackup(),
				defaultPodVolumeBackup().Namespace("ns-1").Name("pvb-2").PodVolumeBackup(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				client          = fake.NewSimpleClientset()
				sharedInformers = informers.NewSharedInformerFactory(client, 0)
				pluginManager   = &pluginmocks.Manager{}
				backupStores    = make(map[string]*persistencemocks.BackupStore)
			)

			c := NewBackupSyncController(
				client.VeleroV1(),
				client.VeleroV1(),
				client.VeleroV1(),
				sharedInformers.Velero().V1().Backups(),
				sharedInformers.Velero().V1().BackupStorageLocations(),
				sharedInformers.Velero().V1().PodVolumeBackups(),
				time.Duration(0),
				test.namespace,
				"",
				func(logrus.FieldLogger) clientmgmt.Manager { return pluginManager },
				velerotest.NewLogger(),
			).(*backupSyncController)

			c.newBackupStore = func(loc *velerov1api.BackupStorageLocation, _ persistence.ObjectStoreGetter, _ logrus.FieldLogger) (persistence.BackupStore, error) {
				// this gets populated just below, prior to exercising the method under test
				return backupStores[loc.Name], nil
			}

			pluginManager.On("CleanupClients").Return(nil)

			for _, location := range test.locations {
				require.NoError(t, sharedInformers.Velero().V1().BackupStorageLocations().Informer().GetStore().Add(location))
				backupStores[location.Name] = &persistencemocks.BackupStore{}
			}

			for _, location := range test.locations {
				backupStore, ok := backupStores[location.Name]
				require.True(t, ok, "no mock backup store for location %s", location.Name)

				backupStore.On("GetRevision").Return("foo", nil)

				var backupNames []string
				for _, bucket := range test.cloudBuckets[location.Spec.ObjectStorage.Bucket] {
					backupNames = append(backupNames, bucket.backup.Name)
					backupStore.On("GetBackupMetadata", bucket.backup.Name).Return(bucket.backup, nil)
					backupStore.On("GetPodVolumeBackups", bucket.backup.Name).Return(bucket.podVolumeBackups, nil)
				}
				backupStore.On("ListBackups").Return(backupNames, nil)
			}

			for _, existingBackup := range test.existingBackups {
				require.NoError(t, sharedInformers.Velero().V1().Backups().Informer().GetStore().Add(existingBackup))

				_, err := client.VeleroV1().Backups(test.namespace).Create(existingBackup)
				require.NoError(t, err)
			}

			for _, existingPodVolumeBackup := range test.existingPodVolumeBackups {
				require.NoError(t, sharedInformers.Velero().V1().PodVolumeBackups().Informer().GetStore().Add(existingPodVolumeBackup))

				_, err := client.VeleroV1().PodVolumeBackups(test.namespace).Create(existingPodVolumeBackup)
				require.NoError(t, err)
			}
			client.ClearActions()

			c.run()

			for bucket, backupDataSet := range test.cloudBuckets {
				// figure out which location this bucket is for; we need this for verification
				// purposes later
				var location *velerov1api.BackupStorageLocation
				for _, loc := range test.locations {
					if loc.Spec.ObjectStorage.Bucket == bucket {
						location = loc
						break
					}
				}
				require.NotNil(t, location)

				// process the cloud backups
				for _, cloudBackupData := range backupDataSet {
					obj, err := client.VeleroV1().Backups(test.namespace).Get(cloudBackupData.backup.Name, metav1.GetOptions{})
					require.NoError(t, err)

					// did this cloud backup already exist in the cluster?
					var existing *velerov1api.Backup
					for _, obj := range test.existingBackups {
						if obj.Name == cloudBackupData.backup.Name {
							existing = obj
							break
						}
					}

					if existing != nil {
						// if this cloud backup already exists in the cluster, make sure that what we get from the
						// client is the existing backup, not the cloud one.

						// verify that the in-cluster backup has its storage location populated, if it's not already.
						expected := existing.DeepCopy()
						expected.Spec.StorageLocation = location.Name

						assert.Equal(t, expected, obj)
					} else {
						// verify that the storage location field and label are set properly
						assert.Equal(t, location.Name, obj.Spec.StorageLocation)

						locationName := location.Name
						if test.longLocationNameEnabled {
							locationName = label.GetValidName(locationName)
						}
						assert.Equal(t, locationName, obj.Labels[velerov1api.StorageLocationLabel])
						assert.Equal(t, true, len(obj.Labels[velerov1api.StorageLocationLabel]) <= validation.DNS1035LabelMaxLength)
					}

					// process the cloud pod volume backups for this backup, if any
					for _, podVolumeBackup := range cloudBackupData.podVolumeBackups {
						objPodVolumeBackup, err := client.VeleroV1().PodVolumeBackups(test.namespace).Get(podVolumeBackup.Name, metav1.GetOptions{})
						require.NoError(t, err)

						// did this cloud pod volume backup already exist in the cluster?
						var existingPodVolumeBackup *velerov1api.PodVolumeBackup
						for _, objPodVolumeBackup := range test.existingPodVolumeBackups {
							if objPodVolumeBackup.Name == podVolumeBackup.Name {
								existingPodVolumeBackup = objPodVolumeBackup
								break
							}
						}

						if existingPodVolumeBackup != nil {
							// if this cloud pod volume backup already exists in the cluster, make sure that what we get from the
							// client is the existing backup, not the cloud one.
							expected := existingPodVolumeBackup.DeepCopy()
							assert.Equal(t, expected, objPodVolumeBackup)
						}
					}
				}
			}
		})
	}
}

func TestDeleteOrphanedBackups(t *testing.T) {
	tests := []struct {
		name            string
		cloudBackups    sets.String
		k8sBackups      []*velerov1api.Backup
		namespace       string
		expectedDeletes sets.String
	}{
		{
			name:         "no overlapping backups",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backupA").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backupB").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backupC").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
			},
			expectedDeletes: sets.NewString("backupA", "backupB", "backupC"),
		},
		{
			name:         "some overlapping backups",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-2").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-C").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
			},
			expectedDeletes: sets.NewString("backup-C"),
		},
		{
			name:         "all overlapping backups",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-2").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-3").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
			},
			expectedDeletes: sets.NewString(),
		},
		{
			name:         "no overlapping backups but including backups that are not complete",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backupA").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("Deleting").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseDeleting).Backup(),
				defaultBackup().Namespace("ns-1").Name("Failed").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseFailed).Backup(),
				defaultBackup().Namespace("ns-1").Name("FailedValidation").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseFailedValidation).Backup(),
				defaultBackup().Namespace("ns-1").Name("InProgress").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseInProgress).Backup(),
				defaultBackup().Namespace("ns-1").Name("New").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseNew).Backup(),
			},
			expectedDeletes: sets.NewString("backupA"),
		},
		{
			name:         "all overlapping backups and all backups that are not complete",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseFailed).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-2").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseFailedValidation).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-3").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseInProgress).Backup(),
			},
			expectedDeletes: sets.NewString(),
		},
		{
			name:         "no completed backups in other locations are deleted",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backup-1").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-2").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-C").Labels(velerov1api.StorageLocationLabel, "default").Phase(velerov1api.BackupPhaseCompleted).Backup(),

				defaultBackup().Namespace("ns-1").Name("backup-4").Labels(velerov1api.StorageLocationLabel, "alternate").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-5").Labels(velerov1api.StorageLocationLabel, "alternate").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-6").Labels(velerov1api.StorageLocationLabel, "alternate").Phase(velerov1api.BackupPhaseCompleted).Backup(),
			},
			expectedDeletes: sets.NewString("backup-C"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				client          = fake.NewSimpleClientset()
				sharedInformers = informers.NewSharedInformerFactory(client, 0)
			)

			c := NewBackupSyncController(
				client.VeleroV1(),
				client.VeleroV1(),
				client.VeleroV1(),
				sharedInformers.Velero().V1().Backups(),
				sharedInformers.Velero().V1().BackupStorageLocations(),
				sharedInformers.Velero().V1().PodVolumeBackups(),
				time.Duration(0),
				test.namespace,
				"",
				nil, // new plugin manager func
				velerotest.NewLogger(),
			).(*backupSyncController)

			expectedDeleteActions := make([]core.Action, 0)

			for _, backup := range test.k8sBackups {
				// add test backup to informer
				require.NoError(t, sharedInformers.Velero().V1().Backups().Informer().GetStore().Add(backup), "Error adding backup to informer")

				// add test backup to client
				_, err := client.VeleroV1().Backups(test.namespace).Create(backup)
				require.NoError(t, err, "Error adding backup to clientset")

				// if we expect this backup to be deleted, set up the expected DeleteAction
				if test.expectedDeletes.Has(backup.Name) {
					actionDelete := core.NewDeleteAction(
						velerov1api.SchemeGroupVersion.WithResource("backups"),
						test.namespace,
						backup.Name,
					)
					expectedDeleteActions = append(expectedDeleteActions, actionDelete)
				}
			}

			c.deleteOrphanedBackups("default", test.cloudBackups, velerotest.NewLogger())

			numBackups, err := numBackups(t, client, c.namespace)
			assert.NoError(t, err)

			expected := len(test.k8sBackups) - len(test.expectedDeletes)
			assert.Equal(t, expected, numBackups)

			velerotest.CompareActions(t, expectedDeleteActions, getDeleteActions(client.Actions()))
		})
	}
}

func TestStorageLabelsInDeleteOrphanedBackups(t *testing.T) {
	longLabelName := "the-really-long-location-name-that-is-much-more-than-63-characters"
	tests := []struct {
		name            string
		cloudBackups    sets.String
		k8sBackups      []*velerov1api.Backup
		namespace       string
		expectedDeletes sets.String
	}{
		{
			name:         "some overlapping backups",
			namespace:    "ns-1",
			cloudBackups: sets.NewString("backup-1", "backup-2", "backup-3"),
			k8sBackups: []*velerov1api.Backup{
				defaultBackup().Namespace("ns-1").Name("backup-1").
					Labels(velerov1api.StorageLocationLabel, "the-really-long-location-name-that-is-much-more-than-63-c69e779").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-2").
					Labels(velerov1api.StorageLocationLabel, "the-really-long-location-name-that-is-much-more-than-63-c69e779").Phase(velerov1api.BackupPhaseCompleted).Backup(),
				defaultBackup().Namespace("ns-1").Name("backup-C").
					Labels(velerov1api.StorageLocationLabel, "the-really-long-location-name-that-is-much-more-than-63-c69e779").Phase(velerov1api.BackupPhaseCompleted).Backup(),
			},
			expectedDeletes: sets.NewString("backup-C"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				client          = fake.NewSimpleClientset()
				sharedInformers = informers.NewSharedInformerFactory(client, 0)
			)

			c := NewBackupSyncController(
				client.VeleroV1(),
				client.VeleroV1(),
				client.VeleroV1(),
				sharedInformers.Velero().V1().Backups(),
				sharedInformers.Velero().V1().BackupStorageLocations(),
				sharedInformers.Velero().V1().PodVolumeBackups(),
				time.Duration(0),
				test.namespace,
				"",
				nil, // new plugin manager func
				velerotest.NewLogger(),
			).(*backupSyncController)

			expectedDeleteActions := make([]core.Action, 0)

			for _, backup := range test.k8sBackups {
				// add test backup to informer
				require.NoError(t, sharedInformers.Velero().V1().Backups().Informer().GetStore().Add(backup), "Error adding backup to informer")

				// add test backup to client
				_, err := client.VeleroV1().Backups(test.namespace).Create(backup)
				require.NoError(t, err, "Error adding backup to clientset")

				// if we expect this backup to be deleted, set up the expected DeleteAction
				if test.expectedDeletes.Has(backup.Name) {
					actionDelete := core.NewDeleteAction(
						velerov1api.SchemeGroupVersion.WithResource("backups"),
						test.namespace,
						backup.Name,
					)
					expectedDeleteActions = append(expectedDeleteActions, actionDelete)
				}
			}

			c.deleteOrphanedBackups(longLabelName, test.cloudBackups, velerotest.NewLogger())

			numBackups, err := numBackups(t, client, c.namespace)
			assert.NoError(t, err)

			expected := len(test.k8sBackups) - len(test.expectedDeletes)
			assert.Equal(t, expected, numBackups)

			velerotest.CompareActions(t, expectedDeleteActions, getDeleteActions(client.Actions()))
		})
	}
}

func TestShouldSync(t *testing.T) {
	c := clock.NewFakeClock(time.Now())

	tests := []struct {
		name                string
		location            *velerov1api.BackupStorageLocation
		backupStoreRevision string
		now                 time.Time
		expectSync          bool
		expectedRevision    string
	}{
		{
			name:                "BSL with no last-synced metadata should sync",
			location:            &velerov1api.BackupStorageLocation{},
			backupStoreRevision: "foo",
			now:                 c.Now(),
			expectSync:          true,
			expectedRevision:    "foo",
		},
		{
			name: "BSL with unchanged revision last synced more than an hour ago should sync",
			location: &velerov1api.BackupStorageLocation{
				Status: velerov1api.BackupStorageLocationStatus{
					LastSyncedRevision: types.UID("foo"),
					LastSyncedTime:     metav1.Time{Time: c.Now().Add(-61 * time.Minute)},
				},
			},
			backupStoreRevision: "foo",
			now:                 c.Now(),
			expectSync:          true,
			expectedRevision:    "foo",
		},
		{
			name: "BSL with unchanged revision last synced less than an hour ago should not sync",
			location: &velerov1api.BackupStorageLocation{
				Status: velerov1api.BackupStorageLocationStatus{
					LastSyncedRevision: types.UID("foo"),
					LastSyncedTime:     metav1.Time{Time: c.Now().Add(-59 * time.Minute)},
				},
			},
			backupStoreRevision: "foo",
			now:                 c.Now(),
			expectSync:          false,
		},
		{
			name: "BSL with different revision than backup store last synced less than an hour ago should sync",
			location: &velerov1api.BackupStorageLocation{
				Status: velerov1api.BackupStorageLocationStatus{
					LastSyncedRevision: types.UID("foo"),
					LastSyncedTime:     metav1.Time{Time: c.Now().Add(-time.Minute)},
				},
			},
			backupStoreRevision: "bar",
			now:                 c.Now(),
			expectSync:          true,
			expectedRevision:    "bar",
		},
		{
			name: "BSL with different revision than backup store last synced more than an hour ago should sync",
			location: &velerov1api.BackupStorageLocation{
				Status: velerov1api.BackupStorageLocationStatus{
					LastSyncedRevision: types.UID("foo"),
					LastSyncedTime:     metav1.Time{Time: c.Now().Add(-61 * time.Minute)},
				},
			},
			backupStoreRevision: "bar",
			now:                 c.Now(),
			expectSync:          true,
			expectedRevision:    "bar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			backupStore := new(persistencemocks.BackupStore)
			if test.backupStoreRevision != "" {
				backupStore.On("GetRevision").Return(test.backupStoreRevision, nil)
			} else {
				backupStore.On("GetRevision").Return("", errors.New("object revision not found"))
			}

			shouldSync, rev := shouldSync(test.location, test.now, backupStore, velerotest.NewLogger())
			assert.Equal(t, test.expectSync, shouldSync)
			assert.Equal(t, test.expectedRevision, rev)
		})
	}

}

func getDeleteActions(actions []core.Action) []core.Action {
	var deleteActions []core.Action
	for _, action := range actions {
		if action.GetVerb() == "delete" {
			deleteActions = append(deleteActions, action)
		}
	}
	return deleteActions
}

func numBackups(t *testing.T, c *fake.Clientset, ns string) (int, error) {
	t.Helper()
	existingK8SBackups, err := c.VeleroV1().Backups(ns).List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(existingK8SBackups.Items), nil
}

func numPodVolumeBackups(t *testing.T, c *fake.Clientset, ns string) (int, error) {
	t.Helper()
	existingK8SPodvolumeBackups, err := c.VeleroV1().PodVolumeBackups(ns).List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(existingK8SPodvolumeBackups.Items), nil
}
