// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HandlerFn is an autogenerated mock type for the HandlerFn type
type HandlerFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: payload
func (_m *HandlerFn) Execute(payload []byte) error {
	ret := _m.Called(payload)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHandlerFn creates a new instance of HandlerFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHandlerFn(t interface {
	mock.TestingT
	Cleanup(func())
}) *HandlerFn {
	mock := &HandlerFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
