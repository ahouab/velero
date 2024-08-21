/*
Copyright The Velero Contributors.

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

package csi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/velero/pkg/features"
	"github.com/vmware-tanzu/velero/pkg/types"
)

func TestCSIFeatureNotEnabledAndPluginIsFromCSI(t *testing.T) {
	features.NewFeatureFlagSet("EnableCSI")
	require.False(t, ShouldSkipAction("abc"))
	require.False(t, ShouldSkipAction("velero.io/csi-pvc-backupper"))

	features.NewFeatureFlagSet("")
	require.True(t, ShouldSkipAction("velero.io/csi-pvc-backupper"))
	require.False(t, ShouldSkipAction("abc"))
}

func TestToSystemAffinity(t *testing.T) {
	tests := []struct {
		name           string
		loadAffinities []*types.LoadAffinity
		expected       *corev1.Affinity
	}{
		{
			name: "loadAffinity is nil",
		},
		{
			name:           "loadAffinity is empty",
			loadAffinities: []*types.LoadAffinity{},
		},
		{
			name: "with match label",
			loadAffinities: []*types.LoadAffinity{
				{
					NodeSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"key-1": "value-1",
						},
					},
				},
			},
			expected: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "key-1",
										Values:   []string{"value-1"},
										Operator: corev1.NodeSelectorOpIn,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "with match expression",
			loadAffinities: []*types.LoadAffinity{
				{
					NodeSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"key-2": "value-2",
						},
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "key-3",
								Values:   []string{"value-3-1", "value-3-2"},
								Operator: metav1.LabelSelectorOpNotIn,
							},
							{
								Key:      "key-4",
								Values:   []string{"value-4-1", "value-4-2", "value-4-3"},
								Operator: metav1.LabelSelectorOpDoesNotExist,
							},
						},
					},
				},
			},
			expected: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "key-2",
										Values:   []string{"value-2"},
										Operator: corev1.NodeSelectorOpIn,
									},
									{
										Key:      "key-3",
										Values:   []string{"value-3-1", "value-3-2"},
										Operator: corev1.NodeSelectorOpNotIn,
									},
									{
										Key:      "key-4",
										Values:   []string{"value-4-1", "value-4-2", "value-4-3"},
										Operator: corev1.NodeSelectorOpDoesNotExist,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			affinity := ToSystemAffinity(test.loadAffinities)
			assert.True(t, reflect.DeepEqual(affinity, test.expected))
		})
	}
}
