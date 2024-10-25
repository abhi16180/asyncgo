# asyncgo

asyncgo is a concurrent task executor designed for simplicity and performance. asyncgo API is similar to
`ExecutorService` in Java. asyncgo abstracts manual handling of goroutines and provides easy to use highly performant
implementations for executing tasks. 

## Features

- **Task Management**: Handles and executes tasks concurrently.
- **Asynchronous processing**: asyncgo provides non-blocking methods to submit and execute any go functions
- **Graceful Shutdown**: Guarantees ongoing tasks are completed before shutting down.
- **Configurable Concurrency**: Allows setting the number of workers to control concurrent task execution.




1. Importing library:
    ```go
   
    package main
    import "github.com/abhi16180/asyncgo"
   
    ```


## Examples

1. Creating executor service and executing tasks

```go
   // without worker pool, executorService.submit(task,args)
   // will spawn new goroutine for each task
   executorService := asyncgo.NewExecutorService()
   future,err:=executorService.submit(...)
   if err !=nil{
	   // handle error
   } else {
	   fmt.Printf("result %v",future.GetResult())
   }   

```
2. Using worker pool

```go

   executorService := asyncgo.NewExecutorService()
   // set worker count and buffer size based on your needs
   workerPool := executorService.NewFixedWorkerPool(&commons.Options{
       WorkerCount: 10,
       BufferSize:  20,
   })
   future,err:=workerPool.Submit(task,arg1,arg2,..argN)
   if err !=nil{
       // handle error
   } else {
       fmt.Printf("result %v",future.GetResult())
   }
```

3. Complete example - worker pool


```go
   
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
       executorService := asyncgo.NewExecutorService()
       // set worker count and buffer size based on your needs
       workerPool := executorService.NewFixedWorkerPool(&commons.Options{
           WorkerCount: 10,
           BufferSize:  20,
       })
       defer workerPool.ShutdownGracefully()
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