package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Worker interface {
}

type WorkerImpl struct {
}

func NewWorker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task, id int64) {
	defer wg.Done()
	log.Println("New worker started")
	for {
		fmt.Println(len(tasks))
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
			fmt.Println(fmt.Sprintf("Buffer is full waiting for current tasks to complete"))
		}
	}
}

/// create a task with result channel
/// pass the task as channel
/// listen to task changes
/// execute tasks
