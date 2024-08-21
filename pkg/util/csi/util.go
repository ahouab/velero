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
	"strings"

	corev1 "k8s.io/api/core/v1"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/features"
	"github.com/vmware-tanzu/velero/pkg/types"
)

const (
	csiPluginNamePrefix = "velero.io/csi-"
)

func ShouldSkipAction(actionName string) bool {
	return !features.IsEnabled(velerov1api.CSIFeatureFlag) && strings.Contains(actionName, csiPluginNamePrefix)
}

func ToSystemAffinity(loadAffinities []*types.LoadAffinity) *corev1.Affinity {
	if len(loadAffinities) == 0 {
		return nil
	}

	result := new(corev1.Affinity)
	result.NodeAffinity = new(corev1.NodeAffinity)
	result.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = new(corev1.NodeSelector)
	result.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms =
		make([]corev1.NodeSelectorTerm, 0)

	for _, loadAffinity := range loadAffinities {
		requirements := []corev1.NodeSelectorRequirement{}
		for k, v := range loadAffinity.NodeSelector.MatchLabels {
			requirements = append(requirements, corev1.NodeSelectorRequirement{
				Key:      k,
				Values:   []string{v},
				Operator: corev1.NodeSelectorOpIn,
			})
		}

		for _, exp := range loadAffinity.NodeSelector.MatchExpressions {
			requirements = append(requirements, corev1.NodeSelectorRequirement{
				Key:      exp.Key,
				Values:   exp.Values,
				Operator: corev1.NodeSelectorOperator(exp.Operator),
			})
		}

		if len(requirements) == 0 {
			return nil
		}

		result.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms =
			append(
				result.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
				corev1.NodeSelectorTerm{
					MatchExpressions: requirements,
				},
			)
	}

	return result
}
