# asyncgo

Asyncgo is zero-dependency asynchronous task executor written in pure go, that prioritises speed and ease of use.

###  Features
- Asynchronous Task Execution: Submit tasks to execute asynchronously and retrieve results.
- No Manual Goroutine Management: Abstracts away the complexity of managing goroutines, and simplifying the code.
- Worker Pool Management: Asyncgo carefully handles worker pool creation & task execution.
- Graceful Shutdown: Guarantees all existing tasks are completed before shutting down the workers
- Task Cancellation: Support for terminating workers.

### Usecases

- Asynchronous HTTP Requests for Microservices
- Background Job Execution
- Infinite concurrent pollling with worker pool (receiving messages from AWS SQS or similar services)


### Documentation

1. Installation
    ```
    go get github.com/abhi16180/asyncgo"
    ```
2. Importing
   ```
   import "github.com/abhi16180/asyncgo"
   ```

### Examples 

1. Executing multiple functions asynchronously 

```go

    package main

    import (
        "github.com/abhi16180/asyncgo"
        "log"
        "time"
    )

    func main() {
        executor := asyncgo.NewExecutor()
        future1 := executor.Submit(func(arg int) (int64, error) {
            time.Sleep(1 * time.Second)
            return int64(arg * arg), nil
        }, 10)
        // first param is function, all remaining params are arguments that needs to be passed for your function
        // if function signature / args do not match, it will result in execution error
        future2 := executor.Submit(func(arg1 int, arg2 int) (int, error) {
            time.Sleep(1 * time.Second)
            return arg1 + arg2, nil
        }, 10, 20)
        // err is execution error, this does not represent error returned by your function
        result1, err := future1.Get()
        if err != nil {
            log.Println(err)
            return
        }
        result2, err := future2.Get()
        if err != nil {
            log.Println(err)
            return
        }
        // result is []interface that contains all the return values including error that is returned by your function
        log.Println(result1, result2)
    }
```

2. Executing large number of tasks with fixed sized worker pool
```go

package main

import (
	"context"
	"github.com/abhi16180/asyncgo"
	"github.com/abhi16180/asyncgo/commons"
	"log"
	"time"
)

func main() {
	executor := asyncgo.NewExecutor()
	workerPool := executor.NewFixedWorkerPool(context.Background(), &commons.Options{
		WorkerCount: 100,
		BufferSize:  100,
	})
        // gracefully terminate all workers
	// guarantees every task is executed
	defer workerPool.Shutdown()
	futures := []*asyncgo.Future{}
	for i := 0; i < 1000; i++ {
		future, err := workerPool.Submit(timeConsumingTask)
		if err != nil {
			log.Println("error while submitting task to worker pool")
			continue
		}
		futures = append(futures, future)
	}

	for _, future := range futures {
		result, err := future.Get()
		if err != nil {
			log.Println("error while getting result from future")
			continue
		}
		log.Println(result)
	}
}

func timeConsumingTask() string {
	time.Sleep(2 * time.Second)
	return "success"
}

```
4. Cancelling worker pool in the middle of execution
```go

package main

import (
	"context"
	"github.com/abhi16180/asyncgo"
	"github.com/abhi16180/asyncgo/commons"
	"log"
	"time"
)

func main() {
	executor := asyncgo.NewExecutor()
	workerPool := executor.NewFixedWorkerPool(context.Background(), &commons.Options{
		WorkerCount: 100,
		BufferSize:  100,
	})

	futures := []*asyncgo.Future{}
	for i := 0; i < 1000; i++ {
		future, err := workerPool.Submit(timeConsumingTask)
		if err != nil {
			log.Println("error while submitting task to worker pool")
			continue
		}
		futures = append(futures, future)
	}
	// terminate worker pool in the middle of task(s) execution
	workerPool.Terminate()
}

func timeConsumingTask() string {
	time.Sleep(2 * time.Second)
	return "success"
}

```


4. For more use-cases and complex examples check out <a href="https://github.com/abhi16180/asyncgo/tree/main/examples">examples</a> section