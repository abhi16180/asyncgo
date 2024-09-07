package wp

import (
	"fmt"
	"log"
)

type ExecutorService interface {
	Submit(function interface{}, args ...interface{}) (*Future, error)
	NewFixedWorkerPool()
}

type ExecutorServiceImpl struct {
}

// Submit Spawns new goroutine everytime this function is called.
// Use workerpool instead if you have huge number of tasks
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
func (e *ExecutorServiceImpl) NewFixedWorkerPool() {

}
func NewExecutorService() ExecutorService {
	return &ExecutorServiceImpl{}
}
