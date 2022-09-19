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
	v1 "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/typed/velero/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeVeleroV1 struct {
	*testing.Fake
}

func (c *FakeVeleroV1) Backups(namespace string) v1.BackupInterface {
	return &FakeBackups{c, namespace}
}

func (c *FakeVeleroV1) BackupRepositories(namespace string) v1.BackupRepositoryInterface {
	return &FakeBackupRepositories{c, namespace}
}

func (c *FakeVeleroV1) BackupStorageLocations(namespace string) v1.BackupStorageLocationInterface {
	return &FakeBackupStorageLocations{c, namespace}
}

func (c *FakeVeleroV1) DeleteBackupRequests(namespace string) v1.DeleteBackupRequestInterface {
	return &FakeDeleteBackupRequests{c, namespace}
}

func (c *FakeVeleroV1) DownloadRequests(namespace string) v1.DownloadRequestInterface {
	return &FakeDownloadRequests{c, namespace}
}

func (c *FakeVeleroV1) PodVolumeBackups(namespace string) v1.PodVolumeBackupInterface {
	return &FakePodVolumeBackups{c, namespace}
}

func (c *FakeVeleroV1) PodVolumeRestores(namespace string) v1.PodVolumeRestoreInterface {
	return &FakePodVolumeRestores{c, namespace}
}

func (c *FakeVeleroV1) Restores(namespace string) v1.RestoreInterface {
	return &FakeRestores{c, namespace}
}

func (c *FakeVeleroV1) Schedules(namespace string) v1.ScheduleInterface {
	return &FakeSchedules{c, namespace}
}

func (c *FakeVeleroV1) ServerStatusRequests(namespace string) v1.ServerStatusRequestInterface {
	return &FakeServerStatusRequests{c, namespace}
}

func (c *FakeVeleroV1) VolumeSnapshotLocations(namespace string) v1.VolumeSnapshotLocationInterface {
	return &FakeVolumeSnapshotLocations{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeVeleroV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
