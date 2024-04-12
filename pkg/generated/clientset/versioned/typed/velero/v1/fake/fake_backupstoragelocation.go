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

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/generated/applyconfiguration/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackupStorageLocations implements BackupStorageLocationInterface
type FakeBackupStorageLocations struct {
	Fake *FakeVeleroV1
	ns   string
}

var backupstoragelocationsResource = v1.SchemeGroupVersion.WithResource("backupstoragelocations")

var backupstoragelocationsKind = v1.SchemeGroupVersion.WithKind("BackupStorageLocation")

// Get takes name of the backupStorageLocation, and returns the corresponding backupStorageLocation object, and an error if there is any.
func (c *FakeBackupStorageLocations) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.BackupStorageLocation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backupstoragelocationsResource, c.ns, name), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// List takes label and field selectors, and returns the list of BackupStorageLocations that match those selectors.
func (c *FakeBackupStorageLocations) List(ctx context.Context, opts metav1.ListOptions) (result *v1.BackupStorageLocationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backupstoragelocationsResource, backupstoragelocationsKind, c.ns, opts), &v1.BackupStorageLocationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.BackupStorageLocationList{ListMeta: obj.(*v1.BackupStorageLocationList).ListMeta}
	for _, item := range obj.(*v1.BackupStorageLocationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backupStorageLocations.
func (c *FakeBackupStorageLocations) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(backupstoragelocationsResource, c.ns, opts))

}

// Create takes the representation of a backupStorageLocation and creates it.  Returns the server's representation of the backupStorageLocation, and an error, if there is any.
func (c *FakeBackupStorageLocations) Create(ctx context.Context, backupStorageLocation *v1.BackupStorageLocation, opts metav1.CreateOptions) (result *v1.BackupStorageLocation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backupstoragelocationsResource, c.ns, backupStorageLocation), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// Update takes the representation of a backupStorageLocation and updates it. Returns the server's representation of the backupStorageLocation, and an error, if there is any.
func (c *FakeBackupStorageLocations) Update(ctx context.Context, backupStorageLocation *v1.BackupStorageLocation, opts metav1.UpdateOptions) (result *v1.BackupStorageLocation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backupstoragelocationsResource, c.ns, backupStorageLocation), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackupStorageLocations) UpdateStatus(ctx context.Context, backupStorageLocation *v1.BackupStorageLocation, opts metav1.UpdateOptions) (*v1.BackupStorageLocation, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backupstoragelocationsResource, "status", c.ns, backupStorageLocation), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// Delete takes name of the backupStorageLocation and deletes it. Returns an error if one occurs.
func (c *FakeBackupStorageLocations) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(backupstoragelocationsResource, c.ns, name, opts), &v1.BackupStorageLocation{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackupStorageLocations) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backupstoragelocationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.BackupStorageLocationList{})
	return err
}

// Patch applies the patch and returns the patched backupStorageLocation.
func (c *FakeBackupStorageLocations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.BackupStorageLocation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backupstoragelocationsResource, c.ns, name, pt, data, subresources...), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied backupStorageLocation.
func (c *FakeBackupStorageLocations) Apply(ctx context.Context, backupStorageLocation *velerov1.BackupStorageLocationApplyConfiguration, opts metav1.ApplyOptions) (result *v1.BackupStorageLocation, err error) {
	if backupStorageLocation == nil {
		return nil, fmt.Errorf("backupStorageLocation provided to Apply must not be nil")
	}
	data, err := json.Marshal(backupStorageLocation)
	if err != nil {
		return nil, err
	}
	name := backupStorageLocation.Name
	if name == nil {
		return nil, fmt.Errorf("backupStorageLocation.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backupstoragelocationsResource, c.ns, *name, types.ApplyPatchType, data), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeBackupStorageLocations) ApplyStatus(ctx context.Context, backupStorageLocation *velerov1.BackupStorageLocationApplyConfiguration, opts metav1.ApplyOptions) (result *v1.BackupStorageLocation, err error) {
	if backupStorageLocation == nil {
		return nil, fmt.Errorf("backupStorageLocation provided to Apply must not be nil")
	}
	data, err := json.Marshal(backupStorageLocation)
	if err != nil {
		return nil, err
	}
	name := backupStorageLocation.Name
	if name == nil {
		return nil, fmt.Errorf("backupStorageLocation.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backupstoragelocationsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1.BackupStorageLocation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.BackupStorageLocation), err
}
