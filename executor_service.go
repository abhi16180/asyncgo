package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
)

const BufferedChannelSize int64 = 20

var wg sync.WaitGroup
var mutex sync.Mutex

//go:generate mockery --name=ExecutorService --output=./mocks --outpkg=mocks
type ExecutorService interface {
	Submit(function interface{}, args ...interface{}) (*Future, error)
	NewFixedWorkerPool(options *Options) WorkerPool
	pushToQueue(task *Task)
}

type Options struct {
	WorkerCount int64
	BufferSize  int64
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

// NewFixedWorkerPool Creates pool of workers with given options. Spawns separate go-routine for queue processor
// *Note* - If you are not sure about bufferSize, do not set it explicitly.
// Default bufferSize will be set to BufferedChannelSize
func (e *ExecutorServiceImpl) NewFixedWorkerPool(options *Options) WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task, options.BufferSize)
	for i := int64(0); i < options.WorkerCount; i++ {
		wg.Add(1)
		go NewWorker(ctx, &wg, taskChan, i)
	}
	wg.Add(1)
	go e.taskQueue.ProcessQueue(options, taskChan)
	return NewWorkerPool(e, taskChan, &wg, cancel)
}

// NewExecutorService Creates new executorService
func NewExecutorService() ExecutorService {
	taskQueue := TaskQueueImpl{}
	return &ExecutorServiceImpl{
		taskQueue: &taskQueue,
	}
}

type WorkerPool struct {
	options  *Options
	executor ExecutorService
	taskChan chan Task
	wg       *sync.WaitGroup
	Cancel   context.CancelFunc
}

func NewWorkerPool(executor *ExecutorServiceImpl, taskChan chan Task, wg *sync.WaitGroup, cancel context.CancelFunc) WorkerPool {
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
