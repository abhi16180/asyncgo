package wp

import (
	"time"
)

// TODO remove hardcoded vals
var channelBuffer = 20

//go:generate mockery --name=TaskQueue --output=./mocks --outpkg=mocks
type TaskQueue interface {
	PushToQueue(task *Task)
	PopTask() *Task
	ProcessQueue(options *Options, taskChannel chan<- Task)
}

type TaskQueueImpl struct {
	size  int
	tasks []Task
}

func (t *TaskQueueImpl) PushToQueue(task *Task) {
	mutex.Lock()
	defer mutex.Unlock()
	t.size++
	t.tasks = append(t.tasks, *task)
}

func (t *TaskQueueImpl) PopTask() *Task {
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

func (t *TaskQueueImpl) ProcessQueue(options *Options, taskChannel chan<- Task) {
	for {
		if int64(len(taskChannel)) >= options.BufferSize {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		task := t.PopTask()
		if task != nil {
			taskChannel <- *task
		}
	}
}
