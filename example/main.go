package main

import (
	"fmt"
	"github.com/abhi16180/asyncgo"
	"log"
	"time"
)

type myStruct struct {
	Value int
}

func main() {
	futures := make([]*asyncgo.Future, 0)
	now := time.Now()
	executor := asyncgo.NewExecutor()
	workerPool := executor.NewFixedWorkerPool(&asyncgo.Options{
		WorkerCount: 10,
		BufferSize:  20,
	})
	defer workerPool.Shutdown()
	for i := 0; i < 10; i++ {
		f, err := workerPool.Submit(testFunction)
		if err != nil {
			log.Println(err)
		} else {
			futures = append(futures, f)
		}
	}
	for i := range futures {
		result := futures[i].GetResult()
		fmt.Println(result)
	}
	fmt.Printf("Time cost %v\n", time.Now().Sub(now))
}

func testFunction() (myStruct, int) {
	time.Sleep(2 * time.Second)
	return myStruct{
		Value: 1,
	}, 10
}
