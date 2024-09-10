package wp

import (
	"time"
)

var channelBuffer = 20

type TaskQueue interface {
	PushToQueue(task *Task)
	PopTask() *Task
	ProcessQueue(taskChannel chan<- Task)
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

func (t *TaskQueueImpl) ProcessQueue(taskChannel chan<- Task) {
	for {
		if len(taskChannel) >= channelBuffer {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		task := t.PopTask()
		if task != nil {
			taskChannel <- *task
		}
	}
}
