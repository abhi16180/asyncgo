package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// TODO add more funcs
var wg sync.WaitGroup
var mutex sync.Mutex

//go:generate mockery --name=ExecutorService --output=./mocks --outpkg=mocks
type ExecutorService interface {
	Submit(function interface{}, args ...interface{}) (*Future, error)
	NewFixedWorkerPool(workers int64) WorkerPool
	pushToQueue(task *Task)
}

type ExecutorServiceImpl struct {
	taskQueue TaskQueue
}

// Submit Spawns new goroutine everytime this function is called.
// If you have large number of tasks use NewFixedWorkerPool instead
func (e *ExecutorServiceImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	task := NewTask(resultChan, function, args)
	go func() {
		err := task.Execute()
		if err != nil {
			log.Default().Println(fmt.Println(fmt.Sprintf("ExecutorService.Submit: execute task err: %v", err)))
		}
	}()
	return NewFuture(resultChan), nil
}

// pushToQueue Adds task to task queue associated with the worker pool
func (t *ExecutorServiceImpl) pushToQueue(task *Task) {
	t.taskQueue.PushToQueue(task)
}

// NewFixedWorkerPool Creates pool of workers with N go-routines. Spawns separate go-routine for queue processor
func (e *ExecutorServiceImpl) NewFixedWorkerPool(workers int64) WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task, 20)
	for i := int64(0); i < workers; i++ {
		wg.Add(1)
		go NewWorker(ctx, &wg, taskChan, i)
	}
	wg.Add(1)
	go e.taskQueue.ProcessQueue(taskChan)
	return NewWorkerPool(e, taskChan, &wg, cancel)
}

func NewExecutorService() ExecutorService {
	taskQueue := TaskQueueImpl{}
	return &ExecutorServiceImpl{
		taskQueue: &taskQueue,
	}
}

type WorkerPool struct {
	executor ExecutorService
	taskChan chan Task
	wg       *sync.WaitGroup
	Cancel   context.CancelFunc
}

func NewWorkerPool(executor ExecutorService, taskChan chan Task, wg *sync.WaitGroup, cancel context.CancelFunc) WorkerPool {
	return WorkerPool{
		executor: executor,
		taskChan: taskChan,
		wg:       wg,
		Cancel:   cancel,
	}
}

// Submit Creates new task from function and adds to task queue. This does not execute the function instantaneously.
// Will be eventually processed by the worker(s). For instantaneous execution, use ExecutorService.Submit
// instead
func (w *WorkerPool) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	task := NewTask(resultChan, function, args)
	w.executor.pushToQueue(&task)
	return NewFuture(resultChan), nil
}
