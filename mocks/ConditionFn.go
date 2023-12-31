// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// ConditionFn is an autogenerated mock type for the ConditionFn type
type ConditionFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *ConditionFn) Execute(_a0 context.Context, _a1 *http.Request) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request) (bool, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *http.Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewConditionFn creates a new instance of ConditionFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConditionFn(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConditionFn {
	mock := &ConditionFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
