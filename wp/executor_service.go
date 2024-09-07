package wp

import (
	"fmt"
	"log"
)

type ExecutorService interface {
	Submit(function interface{}, args ...interface{}) (*Future, error)
	NewFixedWorkerPool(workers int64)
}

type ExecutorServiceImpl struct {
}

// Submit Spawns new goroutine everytime this function is called.
// If you have large number of tasks use NewFixedWorkerPool instead
func (e *ExecutorServiceImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan interface{})
	task := NewTask()
	go func() {
		err := task.Execute(function, args, resultChan)
		if err != nil {
			log.Default().Println(fmt.Println(fmt.Sprintf("ExecutorService.Submit: execute task err: %v", err)))
		}
	}()
	return NewFuture(resultChan), nil
}

// NewFixedWorkerPool WIP
func (e *ExecutorServiceImpl) NewFixedWorkerPool(workers int64) {
	for i := int64(0); i < workers; i++ {
		NewWorker()
	}
}
func NewExecutorService() ExecutorService {
	return &ExecutorServiceImpl{}
}
