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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExecHookApplyConfiguration represents an declarative configuration of the ExecHook type for use
// with apply.
type ExecHookApplyConfiguration struct {
	Container *string           `json:"container,omitempty"`
	Command   []string          `json:"command,omitempty"`
	OnError   *v1.HookErrorMode `json:"onError,omitempty"`
	Timeout   *metav1.Duration  `json:"timeout,omitempty"`
}

// ExecHookApplyConfiguration constructs an declarative configuration of the ExecHook type for use with
// apply.
func ExecHook() *ExecHookApplyConfiguration {
	return &ExecHookApplyConfiguration{}
}

// WithContainer sets the Container field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Container field is set to the value of the last call.
func (b *ExecHookApplyConfiguration) WithContainer(value string) *ExecHookApplyConfiguration {
	b.Container = &value
	return b
}

// WithCommand adds the given value to the Command field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Command field.
func (b *ExecHookApplyConfiguration) WithCommand(values ...string) *ExecHookApplyConfiguration {
	for i := range values {
		b.Command = append(b.Command, values[i])
	}
	return b
}

// WithOnError sets the OnError field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OnError field is set to the value of the last call.
func (b *ExecHookApplyConfiguration) WithOnError(value v1.HookErrorMode) *ExecHookApplyConfiguration {
	b.OnError = &value
	return b
}

// WithTimeout sets the Timeout field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Timeout field is set to the value of the last call.
func (b *ExecHookApplyConfiguration) WithTimeout(value metav1.Duration) *ExecHookApplyConfiguration {
	b.Timeout = &value
	return b
}
