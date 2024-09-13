package wp

import (
	"time"
)

//go:generate mockery --name=taskQueue --output=./mocks --outpkg=mocks
type taskQueue interface {
	// PushToQueue pushes task to taskQueue
	PushToQueue(task *Task)
	// PopTask removes the first item from the queue and returns the pointer to it.
	// If item does not exist, returns nil
	PopTask() *Task
	// ProcessQueue continuously checks the buffered channel's size.
	// If the buffered channel is not full, pops tasks from taskQueue
	// and sends to tasks channel
	ProcessQueue(options *Options, taskChannel chan<- Task)
}

type taskQueueImpl struct {
	size  int
	tasks []Task
}

func (t *taskQueueImpl) PushToQueue(task *Task) {
	mutex.Lock()
	defer mutex.Unlock()
	t.size++
	t.tasks = append(t.tasks, *task)
}

func (t *taskQueueImpl) PopTask() *Task {
	mutex.Lock()
	defer mutex.Unlock()
	if t.size > 0 {
		t.size--
		task := t.tasks[0]
		t.tasks = t.tasks[1:]
		return &task
	}
	return nil
}

func (t *taskQueueImpl) ProcessQueue(options *Options, taskChannel chan<- Task) {
	for {
		if int64(len(taskChannel)) >= options.BufferSize {
			time.Sleep(1 * time.Millisecond)
			continue
		}
		task := t.PopTask()
		if task != nil {
			taskChannel <- *task
		}
	}
}
