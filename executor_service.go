package wp

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
	// Submit Spawns new goroutine everytime this function is called.
	// If you have large number of tasks use NewFixedWorkerPool instead
	Submit(function interface{}, args ...interface{}) (*Future, error)
	// NewFixedWorkerPool Creates pool of workers with given options. Spawns separate go-routine for queue processor
	// *Note* - If you are not sure about bufferSize, do not set it explicitly.
	// Default bufferSize will be set to BufferedChannelSize
	NewFixedWorkerPool(options *Options) WorkerPool
	// pushToQueue Adds task to task queue associated with the worker pool
	pushToQueue(task *Task)
}

type Options struct {
	WorkerCount int64
	BufferSize  int64
}

type ExecutorServiceImpl struct {
	taskQueue taskQueue
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

func (t *ExecutorServiceImpl) pushToQueue(task *Task) {
	t.taskQueue.PushToQueue(task)
}

func (e *ExecutorServiceImpl) NewFixedWorkerPool(options *Options) WorkerPool {
	mutex.Lock()
	defer mutex.Unlock()
	options = GetOrDefaultWorkerPoolOptions(options)
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task, options.BufferSize)
	for i := int64(0); i < options.WorkerCount; i++ {
		wg.Add(1)
		go NewWorker(ctx, &wg, taskChan, i)
	}
	wg.Add(1)
	go e.taskQueue.ProcessQueue(options, taskChan)
	return NewWorkerPool(e, &taskChan, &wg, cancel)
}

// NewExecutorService Creates new executorService
func NewExecutorService() ExecutorService {
	taskQueue := taskQueueImpl{}
	return &ExecutorServiceImpl{
		taskQueue: &taskQueue,
	}
}

func GetOrDefaultWorkerPoolOptions(inputOptions *Options) *Options {
	if inputOptions != nil {
		return inputOptions
	}
	return &Options{
		WorkerCount: int64(runtime.NumCPU()),
		BufferSize:  BufferedChannelSize,
	}
}
