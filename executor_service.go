package quasar

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
)

const BufferedChannelSize int64 = 20

var wg sync.WaitGroup
var mutex sync.Mutex

//go:generate mockery --name=ExecutorService --output=./mocks --outpkg=mocks
type ExecutorService interface {
	// Submit spawns new goroutine everytime this function is called.
	// If you have large number of tasks use NewFixedWorkerPool instead
	Submit(function interface{}, args ...interface{}) (*Future, error)
	// NewFixedWorkerPool creates pool of workers with given options. Spawns separate go-routine for queue processor
	// *Note* - If you are not sure about bufferSize, do not set it explicitly.
	// Default bufferSize will be set to BufferedChannelSize
	NewFixedWorkerPool(options *Options) WorkerPool
	// pushToQueue adds task to task queue associated with the worker pool
}

type Options struct {
	WorkerCount int64
	BufferSize  int64
}

type ExecutorServiceImpl struct {
	shutDown chan interface{}
}

func (e *ExecutorServiceImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	mutex.Lock()
	defer mutex.Unlock()
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

func (e *ExecutorServiceImpl) NewFixedWorkerPool(options *Options) WorkerPool {
	mutex.Lock()
	defer mutex.Unlock()
	options = GetOrDefaultWorkerPoolOptions(options)
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task, options.BufferSize)
	taskQueue := NewTaskQueue()
	wg.Add(1)
	go taskQueue.ProcessQueue(options, taskChan, e.shutDown)
	for i := int64(0); i < options.WorkerCount; i++ {
		wg.Add(1)
		go NewWorker(ctx, &wg, taskChan, i)
	}
	return NewWorkerPool(taskQueue, &taskChan, &wg, cancel, &e.shutDown)
}

// NewExecutorService Creates new executorService
func NewExecutorService() ExecutorService {
	shutDown := make(chan interface{})
	return &ExecutorServiceImpl{
		shutDown: shutDown,
	}
}

func GetOrDefaultWorkerPoolOptions(inputOptions *Options) *Options {
	if inputOptions != nil {
		if inputOptions.WorkerCount == 0 {
			inputOptions.WorkerCount = int64(runtime.NumCPU())
		}
		if inputOptions.BufferSize == 0 {
			inputOptions.BufferSize = BufferedChannelSize
		}
		return inputOptions
	}
	return &Options{
		WorkerCount: int64(runtime.NumCPU()),
		BufferSize:  BufferedChannelSize,
	}
}
