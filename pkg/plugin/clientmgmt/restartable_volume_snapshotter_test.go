/*
Copyright 2018 the Velero contributors.

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

package clientmgmt

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/velero/pkg/plugin/clientmgmt/process"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework/common"
	providermocks "github.com/vmware-tanzu/velero/pkg/plugin/velero/mocks"
)

func TestRestartableGetVolumeSnapshotter(t *testing.T) {
	tests := []struct {
		name          string
		plugin        interface{}
		getError      error
		expectedError string
	}{
		{
			name:          "error getting by kind and name",
			getError:      errors.Errorf("get error"),
			expectedError: "get error",
		},
		{
			name:          "wrong type",
			plugin:        3,
			expectedError: "int is not a VolumeSnapshotter!",
		},
		{
			name:   "happy path",
			plugin: new(providermocks.VolumeSnapshotter),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := new(mockRestartableProcess)
			p.Test(t)
			defer p.AssertExpectations(t)

			name := "aws"
			key := process.KindAndName{Kind: common.PluginKindVolumeSnapshotter, Name: name}
			p.On("GetByKindAndName", key).Return(tc.plugin, tc.getError)

			r := &restartableVolumeSnapshotter{
				key:                 key,
				sharedPluginProcess: p,
			}
			a, err := r.getVolumeSnapshotter()
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.plugin, a)
		})
	}
}

func TestRestartableVolumeSnapshotterReinitialize(t *testing.T) {
	p := new(mockRestartableProcess)
	p.Test(t)
	defer p.AssertExpectations(t)

	name := "aws"
	key := process.KindAndName{Kind: common.PluginKindVolumeSnapshotter, Name: name}
	r := &restartableVolumeSnapshotter{
		key:                 key,
		sharedPluginProcess: p,
		config: map[string]string{
			"color": "blue",
		},
	}

	err := r.Reinitialize(3)
	assert.EqualError(t, err, "int is not a VolumeSnapshotter!")

	volumeSnapshotter := new(providermocks.VolumeSnapshotter)
	volumeSnapshotter.Test(t)
	defer volumeSnapshotter.AssertExpectations(t)

	volumeSnapshotter.On("Init", r.config).Return(errors.Errorf("init error")).Once()
	err = r.Reinitialize(volumeSnapshotter)
	assert.EqualError(t, err, "init error")

	volumeSnapshotter.On("Init", r.config).Return(nil)
	err = r.Reinitialize(volumeSnapshotter)
	assert.NoError(t, err)
}

func TestRestartableVolumeSnapshotterGetDelegate(t *testing.T) {
	p := new(mockRestartableProcess)
	p.Test(t)
	defer p.AssertExpectations(t)

	// Reset error
	p.On("ResetIfNeeded").Return(errors.Errorf("reset error")).Once()
	name := "aws"
	key := process.KindAndName{Kind: common.PluginKindVolumeSnapshotter, Name: name}
	r := &restartableVolumeSnapshotter{
		key:                 key,
		sharedPluginProcess: p,
	}
	a, err := r.getDelegate()
	assert.Nil(t, a)
	assert.EqualError(t, err, "reset error")

	// Happy path
	p.On("ResetIfNeeded").Return(nil)
	volumeSnapshotter := new(providermocks.VolumeSnapshotter)
	volumeSnapshotter.Test(t)
	defer volumeSnapshotter.AssertExpectations(t)
	p.On("GetByKindAndName", key).Return(volumeSnapshotter, nil)

	a, err = r.getDelegate()
	assert.NoError(t, err)
	assert.Equal(t, volumeSnapshotter, a)
}

func TestRestartableVolumeSnapshotterInit(t *testing.T) {
	p := new(mockRestartableProcess)
	p.Test(t)
	defer p.AssertExpectations(t)

	// getVolumeSnapshottererror
	name := "aws"
	key := process.KindAndName{Kind: common.PluginKindVolumeSnapshotter, Name: name}
	r := &restartableVolumeSnapshotter{
		key:                 key,
		sharedPluginProcess: p,
	}
	p.On("GetByKindAndName", key).Return(nil, errors.Errorf("GetByKindAndName error")).Once()

	config := map[string]string{
		"color": "blue",
	}
	err := r.Init(config)
	assert.EqualError(t, err, "GetByKindAndName error")

	// Delegate returns error
	volumeSnapshotter := new(providermocks.VolumeSnapshotter)
	volumeSnapshotter.Test(t)
	defer volumeSnapshotter.AssertExpectations(t)
	p.On("GetByKindAndName", key).Return(volumeSnapshotter, nil)
	volumeSnapshotter.On("Init", config).Return(errors.Errorf("Init error")).Once()

	err = r.Init(config)
	assert.EqualError(t, err, "Init error")

	// wipe this out because the previous failed Init call set it
	r.config = nil

	// Happy path
	volumeSnapshotter.On("Init", config).Return(nil)
	err = r.Init(config)
	assert.NoError(t, err)
	assert.Equal(t, config, r.config)

	// Calling Init twice is forbidden
	err = r.Init(config)
	assert.EqualError(t, err, "already initialized")
}

func TestRestartableVolumeSnapshotterDelegatedFunctions(t *testing.T) {
	pv := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"color": "blue",
		},
	}

	pvToReturn := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"color": "green",
		},
	}

	runRestartableDelegateTests(
		t,
		common.PluginKindVolumeSnapshotter,
		func(key process.KindAndName, p process.RestartableProcess) interface{} {
			return &restartableVolumeSnapshotter{
				key:                 key,
				sharedPluginProcess: p,
			}
		},
		func() mockable {
			return new(providermocks.VolumeSnapshotter)
		},
		restartableDelegateTest{
			function:                "CreateVolumeFromSnapshot",
			inputs:                  []interface{}{"snapshotID", "volumeID", "volumeAZ", to.Int64Ptr(10000)},
			expectedErrorOutputs:    []interface{}{"", errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{"volumeID", errors.Errorf("delegate error")},
		},
		restartableDelegateTest{
			function:                "GetVolumeID",
			inputs:                  []interface{}{pv},
			expectedErrorOutputs:    []interface{}{"", errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{"volumeID", errors.Errorf("delegate error")},
		},
		restartableDelegateTest{
			function:                "SetVolumeID",
			inputs:                  []interface{}{pv, "volumeID"},
			expectedErrorOutputs:    []interface{}{nil, errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{pvToReturn, errors.Errorf("delegate error")},
		},
		restartableDelegateTest{
			function:                "GetVolumeInfo",
			inputs:                  []interface{}{"volumeID", "volumeAZ"},
			expectedErrorOutputs:    []interface{}{"", (*int64)(nil), errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{"volumeType", to.Int64Ptr(10000), errors.Errorf("delegate error")},
		},
		restartableDelegateTest{
			function:                "CreateSnapshot",
			inputs:                  []interface{}{"volumeID", "volumeAZ", map[string]string{"a": "b"}},
			expectedErrorOutputs:    []interface{}{"", errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{"snapshotID", errors.Errorf("delegate error")},
		},
		restartableDelegateTest{
			function:                "DeleteSnapshot",
			inputs:                  []interface{}{"snapshotID"},
			expectedErrorOutputs:    []interface{}{errors.Errorf("reset error")},
			expectedDelegateOutputs: []interface{}{errors.Errorf("delegate error")},
		},
	)
}
