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

package test

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// APIResource stores information about a specific Kubernetes API
// resource.
type APIResource struct {
	Group      string
	Version    string
	Name       string
	ShortName  string
	Namespaced bool
	Items      []metav1.Object
}

// GVR returns a GroupVersionResource representing the resource.
func (r *APIResource) GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    r.Group,
		Version:  r.Version,
		Resource: r.Name,
	}
}

// Pods returns an APIResource describing core/v1's Pods.
func Pods(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "pods",
		ShortName:  "po",
		Namespaced: true,
		Items:      items,
	}
}

func PVCs(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "persistentvolumeclaims",
		ShortName:  "pvc",
		Namespaced: true,
		Items:      items,
	}
}

func PVs(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "persistentvolumes",
		ShortName:  "pv",
		Namespaced: false,
		Items:      items,
	}
}

func Secrets(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "secrets",
		ShortName:  "secrets",
		Namespaced: true,
		Items:      items,
	}
}

func Deployments(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "apps",
		Version:    "v1",
		Name:       "deployments",
		ShortName:  "deploy",
		Namespaced: true,
		Items:      items,
	}
}

func ExtensionsDeployments(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "extensions",
		Version:    "v1",
		Name:       "deployments",
		ShortName:  "deploy",
		Namespaced: true,
		Items:      items,
	}
}

func Namespaces(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "namespaces",
		ShortName:  "ns",
		Namespaced: false,
		Items:      items,
	}
}

func ServiceAccounts(items ...metav1.Object) *APIResource {
	return &APIResource{
		Group:      "",
		Version:    "v1",
		Name:       "serviceaccounts",
		ShortName:  "sa",
		Namespaced: true,
		Items:      items,
	}
}

type ObjectOpts func(metav1.Object)

func NewServiceAccount(ns, name string, opts ...ObjectOpts) *corev1.ServiceAccount {
	obj := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: objectMeta(ns, name),
	}

	for _, opt := range opts {
		opt(obj)
	}

	return obj
}

func NewNamespace(name string, opts ...ObjectOpts) *corev1.Namespace {
	obj := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: objectMeta("", name),
	}

	for _, opt := range opts {
		opt(obj)
	}

	return obj
}

func NewStorageClass(name string, opts ...ObjectOpts) *storagev1.StorageClass {
	obj := &storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage/v1",
		},
		ObjectMeta: objectMeta("", name),
	}

	for _, opt := range opts {
		opt(obj)
	}

	return obj
}

func objectMeta(ns, name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: ns,
		Name:      name,
	}
}

// WithLabels is a functional option that applies the specified
// label keys/values to an object.
func WithLabels(labels ...string) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		objLabels := obj.GetLabels()
		if objLabels == nil {
			objLabels = make(map[string]string)
		}

		if len(labels)%2 != 0 {
			labels = append(labels, "")
		}

		for i := 0; i < len(labels); i += 2 {
			objLabels[labels[i]] = labels[i+1]
		}

		obj.SetLabels(objLabels)
	}
}

// WithAnnotations is a functional option that applies the specified
// annotation keys/values to an object.
func WithAnnotations(vals ...string) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		objAnnotations := obj.GetAnnotations()
		if objAnnotations == nil {
			objAnnotations = make(map[string]string)
		}

		if len(vals)%2 != 0 {
			vals = append(vals, "")
		}

		for i := 0; i < len(vals); i += 2 {
			objAnnotations[vals[i]] = vals[i+1]
		}

		obj.SetAnnotations(objAnnotations)
	}
}

// WithClusterName is a functional option that applies the specified
// cluster name to an object.
func WithClusterName(val string) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		obj.SetClusterName(val)
	}
}

// WithFinalizers is a functional option that applies the specified
// finalizers to an object.
func WithFinalizers(vals ...string) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		obj.SetFinalizers(vals)
	}
}

// WithDeletionTimestamp is a functional option that applies the specified
// deletion timestamp to an object.
func WithDeletionTimestamp(val time.Time) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		obj.SetDeletionTimestamp(&metav1.Time{Time: val})
	}
}

// WithUID is a functional option that applies the specified UID to an object.
func WithUID(val string) func(obj metav1.Object) {
	return func(obj metav1.Object) {
		obj.SetUID(types.UID(val))
	}
}
