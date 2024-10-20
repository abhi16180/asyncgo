package main

import (
	"context"
	"github.com/abhi16180/asyncgo"
	"log"
	"math/rand"
	"time"
)

func main() {
	executor := asyncgo.NewExecutor()
	// set worker count and buffer size according to your needs
	workerPool := executor.NewFixedWorkerPool(context.TODO(), &asyncgo.Options{
		WorkerCount: 10,
		BufferSize:  10,
	})

	defer workerPool.Shutdown()
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 10; i++ {
		_, err := workerPool.Submit(receiveMessage, ctx)
		if err != nil {
			return
		}
	}
	// remove this if you want to run indefinitely
	go stopAfterSometime(cancel)

	// waits until all futures are done executing
	// in this case this will run 10 seconds
	// To run indefinitely just remove stopAfterSometime function
	// You can use this for services like SQS to continuously poll for new messages
	err := workerPool.WaitAll()
	if err != nil {
		log.Println(err)
		return
	}
}

func receiveMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			result := mockProviderService()
			process(result)
		}
	}
}

func mockProviderService() []int {
	time.Sleep(100 * time.Millisecond)
	var result []int
	for i := 0; i < 100; i++ {
		result = append(result, rand.Int())
	}
	return result
}

func process(result []int) {
	sum := 0
	for _, val := range result {
		sum += val
	}
	log.Println(sum)
}

func stopAfterSometime(cancel context.CancelFunc) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		cancel()
	}
}
