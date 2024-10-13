// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	quasar "github.com/abhi16180/quasar"
	mock "github.com/stretchr/testify/mock"
)

// ExecutorService is an autogenerated mock type for the ExecutorService type
type ExecutorService struct {
	mock.Mock
}

// NewFixedWorkerPool provides a mock function with given fields: options
func (_m *ExecutorService) NewFixedWorkerPool(options *quasar.Options) quasar.WorkerPool {
	ret := _m.Called(options)

	if len(ret) == 0 {
		panic("no return value specified for NewFixedWorkerPool")
	}

	var r0 quasar.WorkerPool
	if rf, ok := ret.Get(0).(func(*quasar.Options) quasar.WorkerPool); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(quasar.WorkerPool)
		}
	}

	return r0
}

// Submit provides a mock function with given fields: function, args
func (_m *ExecutorService) Submit(function interface{}, args ...interface{}) (*quasar.Future, error) {
	var _ca []interface{}
	_ca = append(_ca, function)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 *quasar.Future
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) (*quasar.Future, error)); ok {
		return rf(function, args...)
	}
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *quasar.Future); ok {
		r0 = rf(function, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quasar.Future)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(function, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
