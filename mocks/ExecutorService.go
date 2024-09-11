// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	wp "wp"

	mock "github.com/stretchr/testify/mock"
)

// ExecutorService is an autogenerated mock type for the ExecutorService type
type ExecutorService struct {
	mock.Mock
}

// NewFixedWorkerPool provides a mock function with given fields: options
func (_m *ExecutorService) NewFixedWorkerPool(options *wp.Options) wp.WorkerPool {
	ret := _m.Called(options)

	if len(ret) == 0 {
		panic("no return value specified for NewFixedWorkerPool")
	}

	var r0 wp.WorkerPool
	if rf, ok := ret.Get(0).(func(*wp.Options) wp.WorkerPool); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(wp.WorkerPool)
		}
	}

	return r0
}

// Submit provides a mock function with given fields: function, args
func (_m *ExecutorService) Submit(function interface{}, args ...interface{}) (*wp.Future, error) {
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

// pushToQueue provides a mock function with given fields: task
func (_m *ExecutorService) pushToQueue(task *wp.Task) {
	_m.Called(task)
}

// NewExecutorService creates a new instance of ExecutorService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutorService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExecutorService {
	mock := &ExecutorService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
