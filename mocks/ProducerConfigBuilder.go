// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	kafka "github.com/Alibay/go-kit/kafka"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ProducerConfigBuilder is an autogenerated mock type for the ProducerConfigBuilder type
type ProducerConfigBuilder struct {
	mock.Mock
}

// Async provides a mock function with given fields: v
func (_m *ProducerConfigBuilder) Async(v bool) kafka.ProducerConfigBuilder {
	ret := _m.Called(v)

	var r0 kafka.ProducerConfigBuilder
	if rf, ok := ret.Get(0).(func(bool) kafka.ProducerConfigBuilder); ok {
		r0 = rf(v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kafka.ProducerConfigBuilder)
		}
	}

	return r0
}

// BatchSize provides a mock function with given fields: size
func (_m *ProducerConfigBuilder) BatchSize(size int) kafka.ProducerConfigBuilder {
	ret := _m.Called(size)

	var r0 kafka.ProducerConfigBuilder
	if rf, ok := ret.Get(0).(func(int) kafka.ProducerConfigBuilder); ok {
		r0 = rf(size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kafka.ProducerConfigBuilder)
		}
	}

	return r0
}

// BatchTimeout provides a mock function with given fields: to
func (_m *ProducerConfigBuilder) BatchTimeout(to time.Duration) kafka.ProducerConfigBuilder {
	ret := _m.Called(to)

	var r0 kafka.ProducerConfigBuilder
	if rf, ok := ret.Get(0).(func(time.Duration) kafka.ProducerConfigBuilder); ok {
		r0 = rf(to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kafka.ProducerConfigBuilder)
		}
	}

	return r0
}

// Build provides a mock function with given fields:
func (_m *ProducerConfigBuilder) Build() *kafka.ProducerConfig {
	ret := _m.Called()

	var r0 *kafka.ProducerConfig
	if rf, ok := ret.Get(0).(func() *kafka.ProducerConfig); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kafka.ProducerConfig)
		}
	}

	return r0
}

// Retry provides a mock function with given fields: _a0, timeout
func (_m *ProducerConfigBuilder) Retry(_a0 int, timeout time.Duration) kafka.ProducerConfigBuilder {
	ret := _m.Called(_a0, timeout)

	var r0 kafka.ProducerConfigBuilder
	if rf, ok := ret.Get(0).(func(int, time.Duration) kafka.ProducerConfigBuilder); ok {
		r0 = rf(_a0, timeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kafka.ProducerConfigBuilder)
		}
	}

	return r0
}

// NewProducerConfigBuilder creates a new instance of ProducerConfigBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProducerConfigBuilder(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProducerConfigBuilder {
	mock := &ProducerConfigBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
