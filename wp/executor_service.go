package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
)

var wg sync.WaitGroup
var mutex sync.Mutex

type ExecutorService interface {
	Submit(function interface{}, args ...interface{}) (*Future, error)
	NewFixedWorkerPool(workers int64) WorkerPool
}

type ExecutorServiceImpl struct {
}

// Submit Spawns new goroutine everytime this function is called.
// If you have large number of tasks use NewFixedWorkerPool instead
func (e *ExecutorServiceImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan interface{})
	task := NewTask(resultChan, function, args)
	go func() {
		err := task.Execute()
		if err != nil {
			log.Default().Println(fmt.Println(fmt.Sprintf("ExecutorService.Submit: execute task err: %v", err)))
		}
	}()
	return NewFuture(resultChan), nil
}

// NewFixedWorkerPool WIP
func (e *ExecutorServiceImpl) NewFixedWorkerPool(workers int64) WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task, 20)
	for i := int64(0); i < workers; i++ {
		wg.Add(1)
		go NewWorker(ctx, &wg, taskChan, i)
	}
	return NewWorkerPool(e, taskChan, &wg, cancel)
}

func NewExecutorService() ExecutorService {
	return &ExecutorServiceImpl{}
}

type WorkerPool struct {
	executor  ExecutorService
	taskQueue *TaskQueue
	taskChan  chan Task
	wg        *sync.WaitGroup
	Cancel    context.CancelFunc
}

func NewWorkerPool(executor ExecutorService, taskChan chan Task, wg *sync.WaitGroup, cancel context.CancelFunc) WorkerPool {
	return WorkerPool{
		executor: executor,
		taskChan: taskChan,
		wg:       wg,
		Cancel:   cancel,
	}
}

func (w *WorkerPool) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan interface{})
	task := NewTask(resultChan, function, args)
	fmt.Println(fmt.Sprintf("ExecutorService.Submit: submit task: %v", task))
	w.taskChan <- task
	return NewFuture(resultChan), nil
}
