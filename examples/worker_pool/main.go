package main

import (
	"context"
	"github.com/abhi16180/asyncgo"
	"log"
	"time"
)

func main() {
	now := time.Now()
	exe := asyncgo.NewExecutor()
	workerPool := exe.NewFixedWorkerPool(context.Background(), &asyncgo.Options{
		WorkerCount: 50,
		BufferSize:  10,
	})

	var futures []*asyncgo.Future

	for i := 0; i < 100; i++ {
		future, err := workerPool.Submit(someLongTask, i)
		if err != nil {
			// this error is thrown if you call this method after shutting down the worker pool
			log.Printf("error submitting task %d: %v", i, err)
			break
		}
		futures = append(futures, future)
	}
	for _, future := range futures {
		result, err := future.Get()
		if err != nil {
			log.Println("error while executing the function", err)
			continue
		}
		log.Println("result", result)
	}
	// gracefully shutdown all the workers
	workerPool.Shutdown()
	log.Printf("total time taken %v", time.Since(now))
}

func someLongTask(val int) int {
	time.Sleep(2 * time.Second)
	return val
}
