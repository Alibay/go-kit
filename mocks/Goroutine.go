// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	goroutine "github.com/Alibay/go-kit/goroutine"
	logger "github.com/Alibay/go-kit/logger"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Goroutine is an autogenerated mock type for the Goroutine type
type Goroutine struct {
	mock.Mock
}

// Cmp provides a mock function with given fields: component
func (_m *Goroutine) Cmp(component string) goroutine.Goroutine {
	ret := _m.Called(component)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(string) goroutine.Goroutine); ok {
		r0 = rf(component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// Go provides a mock function with given fields: ctx, f
func (_m *Goroutine) Go(ctx context.Context, f func()) {
	_m.Called(ctx, f)
}

// Mth provides a mock function with given fields: method
func (_m *Goroutine) Mth(method string) goroutine.Goroutine {
	ret := _m.Called(method)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(string) goroutine.Goroutine); ok {
		r0 = rf(method)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithLogger provides a mock function with given fields: _a0
func (_m *Goroutine) WithLogger(_a0 logger.CLogger) goroutine.Goroutine {
	ret := _m.Called(_a0)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(logger.CLogger) goroutine.Goroutine); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithLoggerFn provides a mock function with given fields: loggerFn
func (_m *Goroutine) WithLoggerFn(loggerFn logger.CLoggerFunc) goroutine.Goroutine {
	ret := _m.Called(loggerFn)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(logger.CLoggerFunc) goroutine.Goroutine); ok {
		r0 = rf(loggerFn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithRetry provides a mock function with given fields: retry
func (_m *Goroutine) WithRetry(retry int) goroutine.Goroutine {
	ret := _m.Called(retry)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(int) goroutine.Goroutine); ok {
		r0 = rf(retry)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithRetryDelay provides a mock function with given fields: delay
func (_m *Goroutine) WithRetryDelay(delay time.Duration) goroutine.Goroutine {
	ret := _m.Called(delay)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(time.Duration) goroutine.Goroutine); ok {
		r0 = rf(delay)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// NewGoroutine creates a new instance of Goroutine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGoroutine(t interface {
	mock.TestingT
	Cleanup(func())
}) *Goroutine {
	mock := &Goroutine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}