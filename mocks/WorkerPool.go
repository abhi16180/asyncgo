// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	wp "github.com/abhi16180/quasar"

	mock "github.com/stretchr/testify/mock"
)

// WorkerPool is an autogenerated mock type for the WorkerPool type
type WorkerPool struct {
	mock.Mock
}

// GetChannelBufferSize provides a mock function with given fields:
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

// GetPoolSize provides a mock function with given fields:
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

// Submit provides a mock function with given fields: function, args
func (_m *WorkerPool) Submit(function interface{}, args ...interface{}) (*wp.Future, error) {
	var _ca []interface{}
	_ca = append(_ca, function)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 *wp.Future
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*wp.Future, error)); ok {
		return rf(function, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *wp.Future); ok {
		r0 = rf(function, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wp.Future)
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
