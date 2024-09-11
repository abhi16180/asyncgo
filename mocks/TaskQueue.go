// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	wp "wp"

	mock "github.com/stretchr/testify/mock"
)

// TaskQueue is an autogenerated mock type for the TaskQueue type
type TaskQueue struct {
	mock.Mock
}

// PopTask provides a mock function with given fields:
func (_m *TaskQueue) PopTask() *wp.Task {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PopTask")
	}

	var r0 *wp.Task
	if rf, ok := ret.Get(0).(func() *wp.Task); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wp.Task)
		}
	}

	return r0
}

// ProcessQueue provides a mock function with given fields: options, taskChannel
func (_m *TaskQueue) ProcessQueue(options *wp.Options, taskChannel chan<- wp.Task) {
	_m.Called(options, taskChannel)
}

// PushToQueue provides a mock function with given fields: task
func (_m *TaskQueue) PushToQueue(task *wp.Task) {
	_m.Called(task)
}

// NewTaskQueue creates a new instance of TaskQueue. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskQueue(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskQueue {
	mock := &TaskQueue{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
