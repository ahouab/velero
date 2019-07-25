/*
Copyright 2019 the Velero contributors.

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

package restore

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1api "k8s.io/api/core/v1"
	storagev1api "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/heptio/velero/pkg/plugin/velero"
	"github.com/heptio/velero/pkg/test"
)

// TestChangeStorageClassActionExecute runs the ChangeStorageClassAction's Execute
// method and validates that the item's storage class is modified (or not) as expected.
// Validation is done by comparing the result of the Execute method to the test case's
// desired result.
func TestChangeStorageClassActionExecute(t *testing.T) {
	tests := []struct {
		name         string
		pvOrPVC      interface{}
		configMap    *corev1api.ConfigMap
		storageClass *storagev1api.StorageClass
		want         interface{}
		wantErr      error
	}{
		{
			name:    "a valid mapping for a persistent volume is applied correctly",
			pvOrPVC: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "storageclass-2"),
			),
			storageClass: test.NewStorageClass("storageclass-2"),
			want:         test.NewPV("pv-1", test.WithStorageClassName("storageclass-2")),
		},
		{
			name:    "a valid mapping for a persistent volume claim is applied correctly",
			pvOrPVC: test.NewPVC("velero", "pvc-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "storageclass-2"),
			),
			storageClass: test.NewStorageClass("storageclass-2"),
			want:         test.NewPVC("velero", "pvc-1", test.WithStorageClassName("storageclass-2")),
		},
		{
			name:    "when no config map exists for the plugin, the item is returned as-is",
			pvOrPVC: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/some-other-plugin", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "storageclass-2"),
			),
			want: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
		},
		{
			name:    "when no storage class mappings exist in the plugin config map, the item is returned as-is",
			pvOrPVC: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
			),
			want: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
		},
		{
			name:    "when persistent volume has no storage class, the item is returned as-is",
			pvOrPVC: test.NewPV("pv-1"),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "storageclass-2"),
			),
			want: test.NewPV("pv-1"),
		},
		{
			name:    "when persistent volume claim has no storage class, the item is returned as-is",
			pvOrPVC: test.NewPVC("velero", "pvc-1"),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "storageclass-2"),
			),
			want: test.NewPVC("velero", "pvc-1"),
		},
		{
			name:    "when persistent volume's storage class has no mapping in the config map, the item is returned as-is",
			pvOrPVC: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-3", "storageclass-4"),
			),
			want: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
		},
		{
			name:    "when persistent volume claim's storage class has no mapping in the config map, the item is returned as-is",
			pvOrPVC: test.NewPVC("velero", "pvc-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-3", "storageclass-4"),
			),
			want: test.NewPVC("velero", "pvc-1", test.WithStorageClassName("storageclass-1")),
		},
		{
			name:    "when persistent volume's storage class is mapped to a nonexistent storage class, an error is returned",
			pvOrPVC: test.NewPV("pv-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "nonexistent-storage-class"),
			),
			wantErr: errors.New("error getting storage class nonexistent-storage-class from API: storageclasses.storage.k8s.io \"nonexistent-storage-class\" not found"),
		},
		{
			name:    "when persistent volume claim's storage class is mapped to a nonexistent storage class, an error is returned",
			pvOrPVC: test.NewPVC("velero", "pvc-1", test.WithStorageClassName("storageclass-1")),
			configMap: test.NewConfigMap("velero", "change-storage-class",
				test.WithLabels("velero.io/plugin-config", "true", "velero.io/change-storage-class", "RestoreItemAction"),
				test.WithConfigMapData("storageclass-1", "nonexistent-storage-class"),
			),
			wantErr: errors.New("error getting storage class nonexistent-storage-class from API: storageclasses.storage.k8s.io \"nonexistent-storage-class\" not found"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			a := NewChangeStorageClassAction(
				logrus.StandardLogger(),
				clientset.CoreV1().ConfigMaps("velero"),
				clientset.StorageV1().StorageClasses(),
			)

			// set up test data
			if tc.configMap != nil {
				_, err := clientset.CoreV1().ConfigMaps(tc.configMap.Namespace).Create(tc.configMap)
				require.NoError(t, err)
			}

			if tc.storageClass != nil {
				_, err := clientset.StorageV1().StorageClasses().Create(tc.storageClass)
				require.NoError(t, err)
			}

			unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(tc.pvOrPVC)
			require.NoError(t, err)

			input := &velero.RestoreItemActionExecuteInput{
				Item: &unstructured.Unstructured{
					Object: unstructuredMap,
				},
			}

			// execute method under test
			res, err := a.Execute(input)

			// validate for both error and non-error cases
			switch {
			case tc.wantErr != nil:
				assert.EqualError(t, err, tc.wantErr.Error())
			default:
				assert.NoError(t, err)

				wantUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(tc.want)
				require.NoError(t, err)

				assert.Equal(t, &unstructured.Unstructured{Object: wantUnstructured}, res.UpdatedItem)
			}
		})
	}
}
