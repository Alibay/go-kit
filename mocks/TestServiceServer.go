// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "github.com/Alibay/go-kit/grpc"
	mock "github.com/stretchr/testify/mock"
)

// TestServiceServer is an autogenerated mock type for the TestServiceServer type
type TestServiceServer struct {
	mock.Mock
}

// Do provides a mock function with given fields: _a0, _a1
func (_m *TestServiceServer) Do(_a0 context.Context, _a1 *grpc.Empty) (*grpc.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *grpc.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.Empty) (*grpc.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.Empty) *grpc.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithError provides a mock function with given fields: _a0, _a1
func (_m *TestServiceServer) WithError(_a0 context.Context, _a1 *grpc.WithErrorRequest) (*grpc.WithErrorResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *grpc.WithErrorResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.WithErrorRequest) (*grpc.WithErrorResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.WithErrorRequest) *grpc.WithErrorResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.WithErrorResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.WithErrorRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithPanic provides a mock function with given fields: _a0, _a1
func (_m *TestServiceServer) WithPanic(_a0 context.Context, _a1 *grpc.WithPanicRequest) (*grpc.WithPanicResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *grpc.WithPanicResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.WithPanicRequest) (*grpc.WithPanicResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.WithPanicRequest) *grpc.WithPanicResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.WithPanicResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.WithPanicRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedTestServiceServer provides a mock function with given fields:
func (_m *TestServiceServer) mustEmbedUnimplementedTestServiceServer() {
	_m.Called()
}

// NewTestServiceServer creates a new instance of TestServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTestServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *TestServiceServer {
	mock := &TestServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
