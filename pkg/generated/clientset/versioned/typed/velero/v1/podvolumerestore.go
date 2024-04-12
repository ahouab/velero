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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/generated/applyconfiguration/velero/v1"
	scheme "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodVolumeRestoresGetter has a method to return a PodVolumeRestoreInterface.
// A group's client should implement this interface.
type PodVolumeRestoresGetter interface {
	PodVolumeRestores(namespace string) PodVolumeRestoreInterface
}

// PodVolumeRestoreInterface has methods to work with PodVolumeRestore resources.
type PodVolumeRestoreInterface interface {
	Create(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.CreateOptions) (*v1.PodVolumeRestore, error)
	Update(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.UpdateOptions) (*v1.PodVolumeRestore, error)
	UpdateStatus(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.UpdateOptions) (*v1.PodVolumeRestore, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.PodVolumeRestore, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PodVolumeRestoreList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.PodVolumeRestore, err error)
	Apply(ctx context.Context, podVolumeRestore *velerov1.PodVolumeRestoreApplyConfiguration, opts metav1.ApplyOptions) (result *v1.PodVolumeRestore, err error)
	ApplyStatus(ctx context.Context, podVolumeRestore *velerov1.PodVolumeRestoreApplyConfiguration, opts metav1.ApplyOptions) (result *v1.PodVolumeRestore, err error)
	PodVolumeRestoreExpansion
}

// podVolumeRestores implements PodVolumeRestoreInterface
type podVolumeRestores struct {
	client rest.Interface
	ns     string
}

// newPodVolumeRestores returns a PodVolumeRestores
func newPodVolumeRestores(c *VeleroV1Client, namespace string) *podVolumeRestores {
	return &podVolumeRestores{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podVolumeRestore, and returns the corresponding podVolumeRestore object, and an error if there is any.
func (c *podVolumeRestores) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.PodVolumeRestore, err error) {
	result = &v1.PodVolumeRestore{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodVolumeRestores that match those selectors.
func (c *podVolumeRestores) List(ctx context.Context, opts metav1.ListOptions) (result *v1.PodVolumeRestoreList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.PodVolumeRestoreList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podvolumerestores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podVolumeRestores.
func (c *podVolumeRestores) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podvolumerestores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a podVolumeRestore and creates it.  Returns the server's representation of the podVolumeRestore, and an error, if there is any.
func (c *podVolumeRestores) Create(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.CreateOptions) (result *v1.PodVolumeRestore, err error) {
	result = &v1.PodVolumeRestore{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podvolumerestores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podVolumeRestore).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a podVolumeRestore and updates it. Returns the server's representation of the podVolumeRestore, and an error, if there is any.
func (c *podVolumeRestores) Update(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.UpdateOptions) (result *v1.PodVolumeRestore, err error) {
	result = &v1.PodVolumeRestore{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(podVolumeRestore.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podVolumeRestore).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *podVolumeRestores) UpdateStatus(ctx context.Context, podVolumeRestore *v1.PodVolumeRestore, opts metav1.UpdateOptions) (result *v1.PodVolumeRestore, err error) {
	result = &v1.PodVolumeRestore{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(podVolumeRestore.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podVolumeRestore).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the podVolumeRestore and deletes it. Returns an error if one occurs.
func (c *podVolumeRestores) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podVolumeRestores) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podvolumerestores").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched podVolumeRestore.
func (c *podVolumeRestores) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.PodVolumeRestore, err error) {
	result = &v1.PodVolumeRestore{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied podVolumeRestore.
func (c *podVolumeRestores) Apply(ctx context.Context, podVolumeRestore *velerov1.PodVolumeRestoreApplyConfiguration, opts metav1.ApplyOptions) (result *v1.PodVolumeRestore, err error) {
	if podVolumeRestore == nil {
		return nil, fmt.Errorf("podVolumeRestore provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(podVolumeRestore)
	if err != nil {
		return nil, err
	}
	name := podVolumeRestore.Name
	if name == nil {
		return nil, fmt.Errorf("podVolumeRestore.Name must be provided to Apply")
	}
	result = &v1.PodVolumeRestore{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *podVolumeRestores) ApplyStatus(ctx context.Context, podVolumeRestore *velerov1.PodVolumeRestoreApplyConfiguration, opts metav1.ApplyOptions) (result *v1.PodVolumeRestore, err error) {
	if podVolumeRestore == nil {
		return nil, fmt.Errorf("podVolumeRestore provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(podVolumeRestore)
	if err != nil {
		return nil, err
	}

	name := podVolumeRestore.Name
	if name == nil {
		return nil, fmt.Errorf("podVolumeRestore.Name must be provided to Apply")
	}

	result = &v1.PodVolumeRestore{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("podvolumerestores").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
