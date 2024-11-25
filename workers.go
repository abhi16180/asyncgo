package asyncgo

import (
	"context"
	"fmt"
	"github.com/abhi16180/asyncgo/commons"
	"github.com/abhi16180/asyncgo/internal"
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
	// Shutdown guarantees all existing tasks will be executed.
	// No new task(s) will be added to the task queue.
	// Trying to Submit new task will return an error
	Shutdown()
	// Terminate terminates all the workers in worker pool - this is not graceful shutdown
	// Any existing task might not run if this method is called in the middle
	Terminate()
	// WaitAll waits until all futures done executing
	WaitAll() error
}

type WorkerPoolService struct {
	options   *commons.Options
	taskChan  *chan internal.Task
	shutDown  *chan interface{}
	futures   []*Future
	wg        *sync.WaitGroup
	Cancel    context.CancelFunc
	taskQueue internal.Queue
}

func NewWorkerPool(taskQueue internal.Queue, taskChan *chan internal.Task, wg *sync.WaitGroup, cancel context.CancelFunc, shutDown *chan interface{}) WorkerPool {
	return &WorkerPoolService{
		taskChan:  taskChan,
		wg:        wg,
		Cancel:    cancel,
		taskQueue: taskQueue,
		shutDown:  shutDown,
		futures:   []*Future{},
	}
}

func (w *WorkerPoolService) Submit(function interface{}, args ...interface{}) (*Future, error) {
	resultChan := make(chan []interface{})
	errChan := make(chan error)
	task := internal.NewTask(resultChan, errChan, function, args)
	err := w.taskQueue.Push(&task)
	if err != nil {
		return nil, err
	}
	f := NewFuture(resultChan, errChan)
	w.futures = append(w.futures, f)
	return f, nil
}

func (w *WorkerPoolService) PoolSize() int64 {
	return w.options.WorkerCount
}

func (w *WorkerPoolService) ChannelBufferSize() int64 {
	return w.options.BufferSize
}

func (w *WorkerPoolService) Shutdown() {
	*w.shutDown <- true
	_ = w.WaitAll()
}

func (w *WorkerPoolService) Terminate() {
	w.Cancel()
}

func (w *WorkerPoolService) WaitAll() error {
	for i := 0; i < len(w.futures); i++ {
		err := w.futures[i].Wait()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//go:generate mockery --name=Worker --output=./mocks --outpkg=mocks
type Worker interface {
}

type WorkerService struct {
}

// worker creates a new worker which processes tasks from tasks channel
func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan internal.Task, id int64) {
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
