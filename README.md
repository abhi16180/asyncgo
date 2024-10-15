# quasar

Quasar is a concurrent task executor designed for simplicity and performance. Quasar API is similar to
`ExecutorService` in Java. Quasar abstracts manual handling of goroutines and provides easy to use highly performant
implementations for executing tasks. 

## Features

- **Task Management**: Handles and executes tasks concurrently.
- **Asynchronous processing**: Quasar provides non-blocking methods to submit and execute any go functions
- **Graceful Shutdown**: Guarantees ongoing tasks are completed before shutting down.
- **Configurable Concurrency**: Allows setting the number of workers to control concurrent task execution.




1. Importing library:
    ```go
   
    package main
    import "https://github.com/abhi16180/quasar.git"
   
    ```


## Example

```go

package main

import (
	"fmt"
	"github.com/abhi16180/quasar"
	"log"
	"time"
)

type myStruct struct {
	Value int
}

func main() {
	futures := make([]*quasar.Future, 0)
	now := time.Now()
	executorService := quasar.NewExecutorService()
	// set worker count and buffer size based on your needs
	workerPool := executorService.NewFixedWorkerPool(&quasar.Options{
		WorkerCount: 10,
		BufferSize:  20,
	})
	
	for i := 0; i < 10; i++ {
		f, err := workerPool.Submit(testFunction,i,i+1)
		if err != nil {
			log.Println(err)
		} else {
			futures = append(futures, f)
		}
	}
	
	for i, _ := range futures {
		result, executionErr := futures[i].GetResult()
		if executionErr != nil {
			fmt.Println(executionErr)
			continue
		}
		fmt.Println(result)
	}
	fmt.Printf("Time cost %v\n", time.Now().Sub(now))
}

func testFunction(a,b int) (myStruct, int) {
	time.Sleep(2 * time.Second)
	return myStruct{
		Value: a+b,
	}, 10
}
```