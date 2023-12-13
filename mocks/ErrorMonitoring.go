// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	monitoring "github.com/Alibay/go-kit/monitoring"
	mock "github.com/stretchr/testify/mock"
)

// ErrorMonitoring is an autogenerated mock type for the ErrorMonitoring type
type ErrorMonitoring struct {
	mock.Mock
}

// BusinessErrorInc provides a mock function with given fields: errCode
func (_m *ErrorMonitoring) BusinessErrorInc(errCode string) {
	_m.Called(errCode)
}

// Error provides a mock function with given fields: err
func (_m *ErrorMonitoring) Error(err error) {
	_m.Called(err)
}

// GetCollector provides a mock function with given fields:
func (_m *ErrorMonitoring) GetCollector() monitoring.MetricsCollector {
	ret := _m.Called()

	var r0 monitoring.MetricsCollector
	if rf, ok := ret.Get(0).(func() monitoring.MetricsCollector); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(monitoring.MetricsCollector)
		}
	}

	return r0
}

// PanicInc provides a mock function with given fields:
func (_m *ErrorMonitoring) PanicInc() {
	_m.Called()
}

// SystemErrorInc provides a mock function with given fields: errCode
func (_m *ErrorMonitoring) SystemErrorInc(errCode string) {
	_m.Called(errCode)
}

// NewErrorMonitoring creates a new instance of ErrorMonitoring. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewErrorMonitoring(t interface {
	mock.TestingT
	Cleanup(func())
}) *ErrorMonitoring {
	mock := &ErrorMonitoring{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}