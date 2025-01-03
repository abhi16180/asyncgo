// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/abhi16180/asyncgo/commons"

	asyncgo "github.com/abhi16180/asyncgo"

	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// NewFixedWorkerPool provides a mock function with given fields: ctx, options
func (_m *Executor) NewFixedWorkerPool(ctx context.Context, options *commons.Options) asyncgo.WorkerPool {
	ret := _m.Called(ctx, options)

	if len(ret) == 0 {
		panic("no return value specified for NewFixedWorkerPool")
	}

	var r0 asyncgo.WorkerPool
	if rf, ok := ret.Get(0).(func(context.Context, *commons.Options) asyncgo.WorkerPool); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(asyncgo.WorkerPool)
		}
	}

	return r0
}

// Submit provides a mock function with given fields: function, args
func (_m *Executor) Submit(function interface{}, args ...interface{}) *asyncgo.Future {
	var _ca []interface{}
	_ca = append(_ca, function)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 *asyncgo.Future
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) *asyncgo.Future); ok {
		r0 = rf(function, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*asyncgo.Future)
		}
	}

	return r0
}

// NewExecutor creates a new instance of Executor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *Executor {
	mock := &Executor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
