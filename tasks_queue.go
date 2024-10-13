package quasar

import (
	"errors"
	"log"
	"time"
)

//go:generate mockery --name=TaskQueue --output=./mocks --outpkg=mocks
type TaskQueue interface {
	// Push pushes task to TaskQueue
	Push(task *Task) error
	// Pop removes the first item from the queue and returns the pointer to it.
	// If item does not exist, returns nil
	Pop() *Task
	// Process continuously checks the buffered channel's size.
	// If the buffered channel is not full, pops tasks from TaskQueue
	// and sends to tasks channel
	Process(options *Options)
}

type TaskQueueService struct {
	size           int
	rejectNewTasks bool
	tasks          []Task
	taskChannel    *chan Task
	shutDown       *chan interface{}
}

func (t *TaskQueueService) Push(task *Task) error {
	mutex.Lock()
	defer mutex.Unlock()
	if t.rejectNewTasks {
		log.Println("cannot add new task after closing worker pool")
		return errors.New("cannot add new task after closing worker pool")
	}
	t.size++
	t.tasks = append(t.tasks, *task)
	return nil
}

func (t *TaskQueueService) Pop() *Task {
	mutex.Lock()
	defer mutex.Unlock()
	if t.size > 0 {
		t.size--
		task := t.tasks[0]
		t.tasks = t.tasks[1:]
		return &task
	}
	if t.rejectNewTasks {
		// if all tasks are completed and new tasks are rejected close the channel
		log.Println("closing all workers")
		close(*t.shutDown)
		close(*t.taskChannel)
	}
	return nil
}

func (t *TaskQueueService) Process(options *Options) {
	for {
		select {
		case _, ok := <-*t.shutDown:
			mutex.Lock()
			if ok {
				log.Printf("shutting down task queue")
				t.rejectNewTasks = true
			}
			mutex.Unlock()
		default:
			if int64(len(*t.taskChannel)) >= options.BufferSize {
				time.Sleep(1 * time.Millisecond)
				continue
			}
			task := t.Pop()
			if task != nil {
				*t.taskChannel <- *task
			}
		}
	}
}

func NewTaskQueue(taskChan *chan Task, shutDown *chan interface{}) TaskQueue {
	return &TaskQueueService{
		taskChannel: taskChan,
		shutDown:    shutDown,
	}
}
