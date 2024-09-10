package wp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worker interface {
}

type WorkerImpl struct {
}

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
