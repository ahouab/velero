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

// Code generated by informer-gen. DO NOT EDIT.

package v2alpha1

import (
	"context"
	time "time"

	velerov2alpha1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v2alpha1"
	versioned "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/vmware-tanzu/velero/pkg/generated/informers/externalversions/internalinterfaces"
	v2alpha1 "github.com/vmware-tanzu/velero/pkg/generated/listers/velero/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DataUploadInformer provides access to a shared informer and lister for
// DataUploads.
type DataUploadInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2alpha1.DataUploadLister
}

type dataUploadInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDataUploadInformer constructs a new informer for DataUpload type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDataUploadInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDataUploadInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDataUploadInformer constructs a new informer for DataUpload type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDataUploadInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VeleroV2alpha1().DataUploads(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VeleroV2alpha1().DataUploads(namespace).Watch(context.TODO(), options)
			},
		},
		&velerov2alpha1.DataUpload{},
		resyncPeriod,
		indexers,
	)
}

func (f *dataUploadInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDataUploadInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *dataUploadInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&velerov2alpha1.DataUpload{}, f.defaultInformer)
}

func (f *dataUploadInformer) Lister() v2alpha1.DataUploadLister {
	return v2alpha1.NewDataUploadLister(f.Informer().GetIndexer())
}
