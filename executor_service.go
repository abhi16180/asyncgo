package asyncgo

import (
	"context"
	"fmt"
	"github.com/abhi16180/asyncgo/commons"
	"github.com/abhi16180/asyncgo/internal"
	"log"
	"runtime"
	"sync"
	"time"
)

const BufferedChannelSize int64 = 20

var wg sync.WaitGroup
var mutex sync.Mutex

//go:generate mockery --name=Executor --output=./mocks --outpkg=mocks
type Executor interface {
	// Submit spawns new goroutine everytime this function is called.
	// If you have large number of tasks use NewFixedWorkerPool instead
	Submit(function interface{}, args ...interface{}) *Future
	// NewFixedWorkerPool creates pool of workers with given options. Spawns separate go-routine for queue processor
	// *Note* - If you are not sure about bufferSize, do not set it explicitly.
	// Default bufferSize will be set to BufferedChannelSize
	NewFixedWorkerPool(ctx context.Context, options *commons.Options) WorkerPool
	// pushToQueue adds task to task queue associated with the worker pool
}

type ExecutorService struct {
}

// NewExecutor Creates new executorService
func NewExecutor() Executor {
	return &ExecutorService{}
}

func (e *ExecutorService) Submit(function interface{}, args ...interface{}) *Future {
	mutex.Lock()
	defer mutex.Unlock()
	resultChan := make(chan []interface{})
	errChan := make(chan error)
	task := internal.NewTask(resultChan, errChan, function, args)
	go func() {
		err := task.Execute()
		if err != nil {
			log.Default().Println(fmt.Println(fmt.Sprintf("Executor.Submit: execute task err: %v", err)))
		}
	}()
	return NewFuture(resultChan, errChan)
}

func (e *ExecutorService) NewFixedWorkerPool(ctx context.Context, options *commons.Options) WorkerPool {
	mutex.Lock()
	defer mutex.Unlock()
	options = GetOrDefaultWorkerPoolOptions(options)
	ctx, cancel := context.WithCancel(ctx)
	taskChan := make(chan internal.Task, options.BufferSize)
	shutDown := make(chan interface{})
	taskQueue := internal.NewTaskQueue(&taskChan, &shutDown)
	wg.Add(1)
	go taskQueue.Process(&wg, options)
	for i := int64(0); i < options.WorkerCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, taskChan, i)
	}
	return NewWorkerPool(taskQueue, &taskChan, &wg, cancel, &shutDown)
}

func GetOrDefaultWorkerPoolOptions(inputOptions *commons.Options) *commons.Options {
	if inputOptions != nil {
		if inputOptions.WorkerCount == 0 {
			inputOptions.WorkerCount = int64(runtime.NumCPU())
		}
		if inputOptions.BufferSize == 0 {
			inputOptions.BufferSize = BufferedChannelSize
		}
		if inputOptions.IdleSleepDuration == 0 {
			inputOptions.IdleSleepDuration = time.Millisecond * 10
		}
		return inputOptions
	}
	return &commons.Options{
		WorkerCount:       int64(runtime.NumCPU()),
		BufferSize:        BufferedChannelSize,
		IdleSleepDuration: time.Millisecond * 10,
	}
}
