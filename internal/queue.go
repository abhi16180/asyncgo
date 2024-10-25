package internal

import (
	"errors"
	"github.com/abhi16180/asyncgo/commons"
	"log"
	"sync"
	"time"
)

//go:generate mockery --name=Queue --output=../mocks --outpkg=mocks
type Queue interface {
	// Push pushes task to Queue
	Push(task *Task) error
	// Pop removes the first item from the queue and returns the pointer to it.
	// If item does not exist, returns nil
	Pop() *Task
	// Process continuously checks the buffered channel's size.
	// If the buffered channel is not full, pops tasks from Queue
	// and sends to tasks channel
	Process(wg *sync.WaitGroup, options *commons.Options)
}

var mutex sync.Mutex

type QueueService struct {
	size                   int
	shutDownSignalReceived bool
	tasks                  []Task
	taskChannel            *chan Task
	shutDown               *chan interface{}
}

func (t *QueueService) Push(task *Task) error {
	mutex.Lock()
	defer mutex.Unlock()
	if t.shutDownSignalReceived {
		log.Println("cannot add new task after closing worker pool")
		return errors.New("cannot add new task after closing worker pool")
	}
	t.size++
	t.tasks = append(t.tasks, *task)
	return nil
}

func (t *QueueService) Pop() *Task {
	mutex.Lock()
	defer mutex.Unlock()
	if t.size > 0 {
		t.size--
		task := t.tasks[0]
		t.tasks = t.tasks[1:]
		return &task
	}
	if t.shutDownSignalReceived {
		// if all tasks are completed and new tasks are rejected close the channel
		log.Println("closing all workers")
		close(*t.shutDown)
		close(*t.taskChannel)
	}
	return nil
}

func (t *QueueService) Process(wg *sync.WaitGroup, options *commons.Options) {
	defer wg.Done()
	for {
		select {
		case _, ok := <-*t.shutDown:
			mutex.Lock()
			if ok {
				log.Printf("shut down signal received - task queue")
				t.shutDownSignalReceived = true
			}
			mutex.Unlock()
		default:
			if int64(len(*t.taskChannel)) >= options.BufferSize {
				continue
			}
			task := t.Pop()
			if task != nil {
				*t.taskChannel <- *task
			} else {
				if t.shutDownSignalReceived {
					log.Println("closing queue")
					return
				} else {
					time.Sleep(options.IdleSleepDuration)
				}
			}
		}
	}
}

func NewTaskQueue(taskChan *chan Task, shutDown *chan interface{}) Queue {
	return &QueueService{
		taskChannel: taskChan,
		shutDown:    shutDown,
	}
}
