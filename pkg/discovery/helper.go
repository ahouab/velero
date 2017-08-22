/*
Copyright 2017 Heptio Inc.

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

package discovery

import (
	"sort"
	"sync"

	kcmdutil "github.com/heptio/ark/third_party/kubernetes/pkg/kubectl/cmd/util"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

// Helper exposes functions for interacting with the Kubernetes discovery
// API.
type Helper interface {
	// Mapper gets a RESTMapper for the current set of resources retrieved
	// from discovery.
	Mapper() meta.RESTMapper

	// Resources gets the current set of resources retrieved from discovery
	// that are backuppable by Ark.
	Resources() []*metav1.APIResourceList

	// Refresh pulls an updated set of Ark-backuppable resources from the
	// discovery API.
	Refresh() error
}

type helper struct {
	discoveryClient discovery.DiscoveryInterface

	// lock guards mapper and resources
	lock      sync.RWMutex
	mapper    meta.RESTMapper
	resources []*metav1.APIResourceList
}

var _ Helper = &helper{}

func NewHelper(discoveryClient discovery.DiscoveryInterface) (Helper, error) {
	h := &helper{
		discoveryClient: discoveryClient,
	}
	if err := h.Refresh(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *helper) Refresh() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	groupResources, err := discovery.GetAPIGroupResources(h.discoveryClient)
	if err != nil {
		return err
	}
	mapper := discovery.NewRESTMapper(groupResources, dynamic.VersionInterfaces)
	shortcutExpander, err := kcmdutil.NewShortcutExpander(mapper, h.discoveryClient)
	if err != nil {
		return err
	}
	h.mapper = shortcutExpander

	preferredResources, err := h.discoveryClient.ServerPreferredResources()
	if err != nil {
		return err
	}

	h.resources = discovery.FilteredBy(
		discovery.ResourcePredicateFunc(func(groupVersion string, r *metav1.APIResource) bool {
			if groupVersion == api.SchemeGroupVersion.String() {
				return false
			}
			return discovery.SupportsAllVerbs{Verbs: []string{"list", "create"}}.Match(groupVersion, r)
		}),
		preferredResources,
	)

	sortResources(h.resources)

	return nil
}

// sortResources sources resources by moving extensions to the end of the slice. The order of all
// the other resources is preserved.
func sortResources(resources []*metav1.APIResourceList) {
	sort.SliceStable(resources, func(i, j int) bool {
		left := resources[i]
		leftGV, _ := schema.ParseGroupVersion(left.GroupVersion)
		// not checking error because it should be impossible to fail to parse data coming from the
		// apiserver
		if leftGV.Group == "extensions" {
			// always sort extensions at the bottom by saying left is "greater"
			return false
		}

		right := resources[j]
		rightGV, _ := schema.ParseGroupVersion(right.GroupVersion)
		// not checking error because it should be impossible to fail to parse data coming from the
		// apiserver
		if rightGV.Group == "extensions" {
			// always sort extensions at the bottom by saying left is "less"
			return true
		}

		return i < j
	})
}

func (h *helper) Mapper() meta.RESTMapper {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.mapper
}

func (h *helper) Resources() []*metav1.APIResourceList {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.resources
}

// ResolveGroupResource uses the RESTMapper to resolve resource to a fully-qualified
// schema.GroupResource. If the RESTMapper is unable to do so, an error is returned instead.
func ResolveGroupResource(mapper meta.RESTMapper, resource string) (schema.GroupResource, error) {
	gvr, err := mapper.ResourceFor(schema.ParseGroupResource(resource).WithVersion(""))
	if err != nil {
		return schema.GroupResource{}, err
	}
	return gvr.GroupResource(), nil
}
