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

// ServerStatusRequestsGetter has a method to return a ServerStatusRequestInterface.
// A group's client should implement this interface.
type ServerStatusRequestsGetter interface {
	ServerStatusRequests(namespace string) ServerStatusRequestInterface
}

// ServerStatusRequestInterface has methods to work with ServerStatusRequest resources.
type ServerStatusRequestInterface interface {
	Create(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.CreateOptions) (*v1.ServerStatusRequest, error)
	Update(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.UpdateOptions) (*v1.ServerStatusRequest, error)
	UpdateStatus(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.UpdateOptions) (*v1.ServerStatusRequest, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServerStatusRequest, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ServerStatusRequestList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServerStatusRequest, err error)
	Apply(ctx context.Context, serverStatusRequest *velerov1.ServerStatusRequestApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ServerStatusRequest, err error)
	ApplyStatus(ctx context.Context, serverStatusRequest *velerov1.ServerStatusRequestApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ServerStatusRequest, err error)
	ServerStatusRequestExpansion
}

// serverStatusRequests implements ServerStatusRequestInterface
type serverStatusRequests struct {
	client rest.Interface
	ns     string
}

// newServerStatusRequests returns a ServerStatusRequests
func newServerStatusRequests(c *VeleroV1Client, namespace string) *serverStatusRequests {
	return &serverStatusRequests{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the serverStatusRequest, and returns the corresponding serverStatusRequest object, and an error if there is any.
func (c *serverStatusRequests) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ServerStatusRequest, err error) {
	result = &v1.ServerStatusRequest{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServerStatusRequests that match those selectors.
func (c *serverStatusRequests) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ServerStatusRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ServerStatusRequestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serverStatusRequests.
func (c *serverStatusRequests) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a serverStatusRequest and creates it.  Returns the server's representation of the serverStatusRequest, and an error, if there is any.
func (c *serverStatusRequests) Create(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.CreateOptions) (result *v1.ServerStatusRequest, err error) {
	result = &v1.ServerStatusRequest{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serverStatusRequest).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a serverStatusRequest and updates it. Returns the server's representation of the serverStatusRequest, and an error, if there is any.
func (c *serverStatusRequests) Update(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.UpdateOptions) (result *v1.ServerStatusRequest, err error) {
	result = &v1.ServerStatusRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(serverStatusRequest.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serverStatusRequest).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *serverStatusRequests) UpdateStatus(ctx context.Context, serverStatusRequest *v1.ServerStatusRequest, opts metav1.UpdateOptions) (result *v1.ServerStatusRequest, err error) {
	result = &v1.ServerStatusRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(serverStatusRequest.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serverStatusRequest).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the serverStatusRequest and deletes it. Returns an error if one occurs.
func (c *serverStatusRequests) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serverStatusRequests) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serverstatusrequests").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched serverStatusRequest.
func (c *serverStatusRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServerStatusRequest, err error) {
	result = &v1.ServerStatusRequest{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied serverStatusRequest.
func (c *serverStatusRequests) Apply(ctx context.Context, serverStatusRequest *velerov1.ServerStatusRequestApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ServerStatusRequest, err error) {
	if serverStatusRequest == nil {
		return nil, fmt.Errorf("serverStatusRequest provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(serverStatusRequest)
	if err != nil {
		return nil, err
	}
	name := serverStatusRequest.Name
	if name == nil {
		return nil, fmt.Errorf("serverStatusRequest.Name must be provided to Apply")
	}
	result = &v1.ServerStatusRequest{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *serverStatusRequests) ApplyStatus(ctx context.Context, serverStatusRequest *velerov1.ServerStatusRequestApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ServerStatusRequest, err error) {
	if serverStatusRequest == nil {
		return nil, fmt.Errorf("serverStatusRequest provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(serverStatusRequest)
	if err != nil {
		return nil, err
	}

	name := serverStatusRequest.Name
	if name == nil {
		return nil, fmt.Errorf("serverStatusRequest.Name must be provided to Apply")
	}

	result = &v1.ServerStatusRequest{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("serverstatusrequests").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
