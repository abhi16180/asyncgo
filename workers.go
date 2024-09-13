package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

//go:generate mockery --name=WorkerPool --output=./mocks --outpkg=mocks
type WorkerPool interface {
	// Submit creates new task from function and adds to task queue. This does not execute the function instantaneously.
	// Will be eventually processed by the worker(s). For instantaneous execution, use ExecutorService.Submit
	// instead
	Submit(function interface{}, args ...interface{}) (*Future, error)
	// GetPoolSize returns the current worker pool size
	GetPoolSize() int64
	// GetChannelBufferSize returns the current channel buffer size
	GetChannelBufferSize() int64
}
type WorkerPoolImpl struct {
	options  *Options
	executor ExecutorService
	taskChan *chan Task
	wg       *sync.WaitGroup
	Cancel   context.CancelFunc
}

func NewWorkerPool(executor *ExecutorServiceImpl, taskChan *chan Task, wg *sync.WaitGroup, cancel context.CancelFunc) WorkerPool {
	return &WorkerPoolImpl{
		executor: executor,
		taskChan: taskChan,
		wg:       wg,
		Cancel:   cancel,
	}
}

func (w *WorkerPoolImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	task := NewTask(resultChan, function, args)
	w.executor.pushToQueue(&task)
	return NewFuture(resultChan), nil
}

func (w *WorkerPoolImpl) GetPoolSize() int64 {
	return w.options.WorkerCount
}

func (w *WorkerPoolImpl) GetChannelBufferSize() int64 {
	return w.options.BufferSize
}

//go:generate mockery --name=Worker --output=./mocks --outpkg=mocks
type Worker interface {
}

type WorkerImpl struct {
}

// NewWorker creates a new worker which processes tasks from tasks channel
func NewWorker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task, id int64) {
	defer wg.Done()
	log.Println("New worker started")
	for {
		select {
		case task := <-tasks:
			log.Println(fmt.Sprintf("Worker %d received task", id))
			if err := task.Execute(); err != nil {
				log.Println(fmt.Sprintf("Worker %d encountered error: %v", id, err))
			}
		case <-ctx.Done():
			log.Println("Worker", id, "exiting - context canceled")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
