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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DataDownloadLister helps list DataDownloads.
// All objects returned here must be treated as read-only.
type DataDownloadLister interface {
	// List lists all DataDownloads in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DataDownload, err error)
	// DataDownloads returns an object that can list and get DataDownloads.
	DataDownloads(namespace string) DataDownloadNamespaceLister
	DataDownloadListerExpansion
}

// dataDownloadLister implements the DataDownloadLister interface.
type dataDownloadLister struct {
	indexer cache.Indexer
}

// NewDataDownloadLister returns a new DataDownloadLister.
func NewDataDownloadLister(indexer cache.Indexer) DataDownloadLister {
	return &dataDownloadLister{indexer: indexer}
}

// List lists all DataDownloads in the indexer.
func (s *dataDownloadLister) List(selector labels.Selector) (ret []*v1alpha1.DataDownload, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DataDownload))
	})
	return ret, err
}

// DataDownloads returns an object that can list and get DataDownloads.
func (s *dataDownloadLister) DataDownloads(namespace string) DataDownloadNamespaceLister {
	return dataDownloadNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DataDownloadNamespaceLister helps list and get DataDownloads.
// All objects returned here must be treated as read-only.
type DataDownloadNamespaceLister interface {
	// List lists all DataDownloads in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DataDownload, err error)
	// Get retrieves the DataDownload from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.DataDownload, error)
	DataDownloadNamespaceListerExpansion
}

// dataDownloadNamespaceLister implements the DataDownloadNamespaceLister
// interface.
type dataDownloadNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DataDownloads in the indexer for a given namespace.
func (s dataDownloadNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DataDownload, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DataDownload))
	})
	return ret, err
}

// Get retrieves the DataDownload from the indexer for a given namespace and name.
func (s dataDownloadNamespaceLister) Get(name string) (*v1alpha1.DataDownload, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("datadownload"), name)
	}
	return obj.(*v1alpha1.DataDownload), nil
}
