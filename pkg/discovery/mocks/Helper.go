// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	schema "k8s.io/apimachinery/pkg/runtime/schema"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	version "k8s.io/apimachinery/pkg/version"
)

// Helper is an autogenerated mock type for the Helper type
type Helper struct {
	mock.Mock
}

// APIGroups provides a mock function with given fields:
func (_m *Helper) APIGroups() []v1.APIGroup {
	ret := _m.Called()

	var r0 []v1.APIGroup
	if rf, ok := ret.Get(0).(func() []v1.APIGroup); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]v1.APIGroup)
		}
	}

	return r0
}

// KindFor provides a mock function with given fields: input
func (_m *Helper) KindFor(input schema.GroupVersionKind) (schema.GroupVersionResource, v1.APIResource, error) {
	ret := _m.Called(input)

	var r0 schema.GroupVersionResource
	var r1 v1.APIResource
	var r2 error
	if rf, ok := ret.Get(0).(func(schema.GroupVersionKind) (schema.GroupVersionResource, v1.APIResource, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(schema.GroupVersionKind) schema.GroupVersionResource); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(schema.GroupVersionResource)
	}

	if rf, ok := ret.Get(1).(func(schema.GroupVersionKind) v1.APIResource); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Get(1).(v1.APIResource)
	}

	if rf, ok := ret.Get(2).(func(schema.GroupVersionKind) error); ok {
		r2 = rf(input)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Refresh provides a mock function with given fields:
func (_m *Helper) Refresh() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResourceFor provides a mock function with given fields: input
func (_m *Helper) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, v1.APIResource, error) {
	ret := _m.Called(input)

	var r0 schema.GroupVersionResource
	var r1 v1.APIResource
	var r2 error
	if rf, ok := ret.Get(0).(func(schema.GroupVersionResource) (schema.GroupVersionResource, v1.APIResource, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(schema.GroupVersionResource) schema.GroupVersionResource); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(schema.GroupVersionResource)
	}

	if rf, ok := ret.Get(1).(func(schema.GroupVersionResource) v1.APIResource); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Get(1).(v1.APIResource)
	}

	if rf, ok := ret.Get(2).(func(schema.GroupVersionResource) error); ok {
		r2 = rf(input)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Resources provides a mock function with given fields:
func (_m *Helper) Resources() []*v1.APIResourceList {
	ret := _m.Called()

	var r0 []*v1.APIResourceList
	if rf, ok := ret.Get(0).(func() []*v1.APIResourceList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.APIResourceList)
		}
	}

	return r0
}

// ServerVersion provides a mock function with given fields:
func (_m *Helper) ServerVersion() *version.Info {
	ret := _m.Called()

	var r0 *version.Info
	if rf, ok := ret.Get(0).(func() *version.Info); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*version.Info)
		}
	}

	return r0
}

// NewHelper creates a new instance of Helper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHelper(t interface {
	mock.TestingT
	Cleanup(func())
}) *Helper {
	mock := &Helper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// serverResourcesInterface is an autogenerated mock type for the serverResourcesInterface type
type serverResourcesInterface struct {
	mock.Mock
}

// ServerGroupsAndResources provides a mock function with given fields:
func (_m *serverResourcesInterface) ServerGroupsAndResources() ([]*v1.APIGroup, []*v1.APIResourceList, error) {
	ret := _m.Called()

	var r0 []*v1.APIGroup
	var r1 []*v1.APIResourceList
	var r2 error
	if rf, ok := ret.Get(0).(func() ([]*v1.APIGroup, []*v1.APIResourceList, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*v1.APIGroup); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.APIGroup)
		}
	}

	if rf, ok := ret.Get(1).(func() []*v1.APIResourceList); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*v1.APIResourceList)
		}
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ServerPreferredResources provides a mock function with given fields:
func (_m *serverResourcesInterface) ServerPreferredResources() ([]*v1.APIResourceList, error) {
	ret := _m.Called()

	var r0 []*v1.APIResourceList
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*v1.APIResourceList, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*v1.APIResourceList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.APIResourceList)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewServerResourcesInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewServerResourcesInterface creates a new instance of serverResourcesInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewServerResourcesInterface(t mockConstructorTestingTnewServerResourcesInterface) *serverResourcesInterface {
	mock := &serverResourcesInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
