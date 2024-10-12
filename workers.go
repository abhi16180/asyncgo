package quasar

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
	// PoolSize returns the current worker pool size
	PoolSize() int64
	// ChannelBufferSize returns the current channel buffer size
	ChannelBufferSize() int64
	// Terminate Terminates all the workers in worker pool
	// TODO gracefully terminate
	Terminate()
}

type WorkerPoolImpl struct {
	options   *Options
	taskChan  *chan Task
	shutDown  *chan interface{}
	wg        *sync.WaitGroup
	Cancel    context.CancelFunc
	taskQueue TaskQueue
}

func NewWorkerPool(taskQueue TaskQueue, taskChan *chan Task, wg *sync.WaitGroup, cancel context.CancelFunc, shutDown *chan interface{}) WorkerPool {
	return &WorkerPoolImpl{
		taskChan:  taskChan,
		wg:        wg,
		Cancel:    cancel,
		taskQueue: taskQueue,
		shutDown:  shutDown,
	}
}

func (w *WorkerPoolImpl) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	task := NewTask(resultChan, function, args)
	w.taskQueue.PushToQueue(&task)
	return NewFuture(resultChan), nil
}

func (w *WorkerPoolImpl) PoolSize() int64 {
	return w.options.WorkerCount
}

func (w *WorkerPoolImpl) ChannelBufferSize() int64 {
	return w.options.BufferSize
}

func (w *WorkerPoolImpl) Terminate() {
	//w.Cancel()
	*w.shutDown <- true
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
