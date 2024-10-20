package asyncgo

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
	// Will be eventually processed by the worker(s). For instantaneous execution, use Executor.Submit
	// instead
	Submit(function interface{}, args ...interface{}) (*Future, error)
	// PoolSize returns the current worker pool size
	PoolSize() int64
	// ChannelBufferSize returns the current channel buffer size
	ChannelBufferSize() int64
	// ShutdownGracefully -  guarantees all existing tasks will be executed.
	// No new task(s) will be added to the task queue.
	// Trying to Submit new task will return an error
	Shutdown()
	// Terminate terminates all the workers in worker pool - this is not graceful shutdown
	// Any existing task might not run if this method is called in the middle
	Terminate()
}

type WorkerPoolService struct {
	options   *Options
	taskChan  *chan Task
	shutDown  *chan interface{}
	wg        *sync.WaitGroup
	Cancel    context.CancelFunc
	taskQueue TaskQueue
}

func NewWorkerPool(taskQueue TaskQueue, taskChan *chan Task, wg *sync.WaitGroup, cancel context.CancelFunc, shutDown *chan interface{}) WorkerPool {
	return &WorkerPoolService{
		taskChan:  taskChan,
		wg:        wg,
		Cancel:    cancel,
		taskQueue: taskQueue,
		shutDown:  shutDown,
	}
}

func (w *WorkerPoolService) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	task := NewTask(resultChan, function, args)
	err := w.taskQueue.Push(&task)
	if err != nil {
		return nil, err
	}
	return NewFuture(resultChan), nil
}

func (w *WorkerPoolService) PoolSize() int64 {
	return w.options.WorkerCount
}

func (w *WorkerPoolService) ChannelBufferSize() int64 {
	return w.options.BufferSize
}

func (w *WorkerPoolService) Shutdown() {
	*w.shutDown <- true
}

func (w *WorkerPoolService) Terminate() {
	w.Cancel()
}

//go:generate mockery --name=Worker --output=./mocks --outpkg=mocks
type Worker interface {
}

type WorkerService struct {
}

// worker creates a new worker which processes tasks from tasks channel
func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task, id int64) {
	log.Printf("worker %v started", id)
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				log.Printf("channel is closed, stopping worker %v", id)
				return
			}
			log.Println(fmt.Sprintf("Worker %d received task", id))
			if err := task.Execute(); err != nil {
				log.Println(fmt.Sprintf("Worker %d encountered error: %v", id, err))
			}
		case <-ctx.Done():
			log.Println("worker", id, "exiting - context canceled")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
