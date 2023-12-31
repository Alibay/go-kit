// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	go_kithttp "github.com/Alibay/go-kit/http"

	mock "github.com/stretchr/testify/mock"
)

// ResourcePolicyManager is an autogenerated mock type for the ResourcePolicyManager type
type ResourcePolicyManager struct {
	mock.Mock
}

// GetRequestedResources provides a mock function with given fields: ctx, routeId, r
func (_m *ResourcePolicyManager) GetRequestedResources(ctx context.Context, routeId string, r *http.Request) ([]*go_kithttp.AuthorizationResource, error) {
	ret := _m.Called(ctx, routeId, r)

	var r0 []*go_kithttp.AuthorizationResource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *http.Request) ([]*go_kithttp.AuthorizationResource, error)); ok {
		return rf(ctx, routeId, r)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *http.Request) []*go_kithttp.AuthorizationResource); ok {
		r0 = rf(ctx, routeId, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*go_kithttp.AuthorizationResource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *http.Request) error); ok {
		r1 = rf(ctx, routeId, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterResourceMapping provides a mock function with given fields: routeId, policies
func (_m *ResourcePolicyManager) RegisterResourceMapping(routeId string, policies ...go_kithttp.ResourcePolicy) {
	_va := make([]interface{}, len(policies))
	for _i := range policies {
		_va[_i] = policies[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, routeId)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// NewResourcePolicyManager creates a new instance of ResourcePolicyManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResourcePolicyManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResourcePolicyManager {
	mock := &ResourcePolicyManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
