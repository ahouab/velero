/*
Copyright the Velero contributors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

// DeleteBackupRequestStatusApplyConfiguration represents an declarative configuration of the DeleteBackupRequestStatus type for use
// with apply.
type DeleteBackupRequestStatusApplyConfiguration struct {
	Phase  *v1.DeleteBackupRequestPhase `json:"phase,omitempty"`
	Errors []string                     `json:"errors,omitempty"`
}

// DeleteBackupRequestStatusApplyConfiguration constructs an declarative configuration of the DeleteBackupRequestStatus type for use with
// apply.
func DeleteBackupRequestStatus() *DeleteBackupRequestStatusApplyConfiguration {
	return &DeleteBackupRequestStatusApplyConfiguration{}
}

// WithPhase sets the Phase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Phase field is set to the value of the last call.
func (b *DeleteBackupRequestStatusApplyConfiguration) WithPhase(value v1.DeleteBackupRequestPhase) *DeleteBackupRequestStatusApplyConfiguration {
	b.Phase = &value
	return b
}

// WithErrors adds the given value to the Errors field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Errors field.
func (b *DeleteBackupRequestStatusApplyConfiguration) WithErrors(values ...string) *DeleteBackupRequestStatusApplyConfiguration {
	for i := range values {
		b.Errors = append(b.Errors, values[i])
	}
	return b
}
