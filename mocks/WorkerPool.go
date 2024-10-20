// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	asyncgo "github.com/abhi16180/asyncgo"
	mock "github.com/stretchr/testify/mock"
)

// WorkerPool is an autogenerated mock type for the WorkerPool type
type WorkerPool struct {
	mock.Mock
}

// ChannelBufferSize provides a mock function with given fields:
func (_m *WorkerPool) ChannelBufferSize() int64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ChannelBufferSize")
	}

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// PoolSize provides a mock function with given fields:
func (_m *WorkerPool) PoolSize() int64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PoolSize")
	}

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *WorkerPool) Shutdown() {
	_m.Called()
}

// Submit provides a mock function with given fields: function, args
func (_m *WorkerPool) Submit(function interface{}, args ...interface{}) (*asyncgo.Future, error) {
	var _ca []interface{}
	_ca = append(_ca, function)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 *asyncgo.Future
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*asyncgo.Future, error)); ok {
		return rf(function, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *asyncgo.Future); ok {
		r0 = rf(function, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*asyncgo.Future)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(function, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Terminate provides a mock function with given fields:
func (_m *WorkerPool) Terminate() {
	_m.Called()
}

// WaitAll provides a mock function with given fields:
func (_m *WorkerPool) WaitAll() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WaitAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWorkerPool creates a new instance of WorkerPool. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWorkerPool(t interface {
	mock.TestingT
	Cleanup(func())
}) *WorkerPool {
	mock := &WorkerPool{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
