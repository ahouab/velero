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

package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestParseSecurityContext(t *testing.T) {
	type args struct {
		runAsUser  string
		runAsGroup string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected *corev1.SecurityContext
	}{
		{"valid security context", args{"1001", "999"}, false, &corev1.SecurityContext{
			RunAsUser:  pointInt64(1001),
			RunAsGroup: pointInt64(999),
		}},
		{"security context without runAsGroup", args{"1001", ""}, false, &corev1.SecurityContext{
			RunAsUser: pointInt64(1001),
		}},
		{"security context without runAsUser", args{"", "999"}, false, &corev1.SecurityContext{
			RunAsGroup: pointInt64(999),
		}},
		{"empty context without runAsUser", args{"", ""}, false, &corev1.SecurityContext{}},
		{"invalid security context runAsUser", args{"not a number", ""}, true, nil},
		{"invalid security context runAsGroup", args{"", "not a number"}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSecurityContext(tt.args.runAsUser, tt.args.runAsGroup)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.expected == nil {
				tt.expected = &corev1.SecurityContext{}
			}

			assert.Equal(t, tt.expected.RunAsUser, got.RunAsUser)
			assert.Equal(t, tt.expected.RunAsGroup, got.RunAsGroup)
		})
	}
}

func pointInt64(i int64) *int64 {
	return &i
}
