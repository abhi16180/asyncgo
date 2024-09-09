package wp

var channelBuffer = 20

type TaskQueue interface {
	PushToQueue(task *Task)
	PopTask() *Task
	ProcessQueue()
}

type TaskQueueImpl struct {
	tasks []Task
}

func (t *TaskQueueImpl) PushToQueue(task *Task) {
	mutex.Lock()
	defer mutex.Unlock()
	t.tasks = append(t.tasks, *task)
}

func (t *TaskQueueImpl) PopTask() *Task {
	task := t.tasks[0]
	t.tasks = t.tasks[1:]
	return &task
}

func (t *TaskQueueImpl) ProcessQueue(taskChannel chan<- Task) {
	mutex.Lock()
	defer mutex.Unlock()
	for i := 0; i < len(t.tasks); i++ {
		if len(taskChannel) < channelBuffer {
			taskChannel <- *t.PopTask()
			continue
		}
	}
}
