package main

import (
	"fmt"
	"github.com/abhi16180/quasar"
	"time"
)

type S struct {
	V int
}

func main() {
	futures := make([]*quasar.Future, 0)
	now := time.Now()
	executorService := quasar.NewExecutorService()
	workerPool := executorService.NewFixedWorkerPool(&quasar.Options{
		WorkerCount: 10,
		BufferSize:  20,
	})

	for i := 0; i < 100; i++ {
		f, _ := workerPool.Submit(testFunction)
		futures = append(futures, f)
	}

	for i := 0; i < 100; i++ {
		result, executionErr := futures[i].GetResult()
		if executionErr != nil {
			fmt.Println(executionErr)
			continue
		}
		fmt.Println(result)
	}
	fmt.Printf("Time cost %v\n", time.Now().Sub(now))
}

func testFunction() (S, int) {
	time.Sleep(2 * time.Second)
	return S{
		V: 1,
	}, 10
}

//type ExecutorService struct {
//}
//
//func NewExecutorService() *ExecutorService {
//	return &ExecutorService{}
//}
//
//func (e *ExecutorService) Submit(fn interface{}, args ...interface{}) (*Future, error) {
//	if reflect.TypeOf(fn).Kind() != reflect.Func {
//		return nil, fmt.Errorf("fn must be a function")
//	}
//	ch := make(chan interface{})
//	if len(args) > 0 {
//		go e.run(ch, fn, args)
//	} else {
//		go e.run(ch, fn)
//	}
//	return NewFuture(ch), nil
//}
//
//func (e *ExecutorService) run(ch chan<- interface{}, fn interface{}, args ...interface{}) {
//	val := reflect.ValueOf(fn)
//	argSlice := make([]reflect.Value, len(args))
//	for i, arg := range args {
//		argSlice[i] = reflect.ValueOf(arg)
//	}
//	if len(argSlice) > 0 {
//		result := val.Call(argSlice)
//		ch <- result
//	}
//	ch <- nil
//}
//
//type Future struct {
//	resultChan <-chan interface{}
//}
//
//func NewFuture(ch <-chan interface{}) *Future {
//	return &Future{
//		resultChan: ch,
//	}
//}
//func (f *Future) Get() interface{} {
//	return <-f.resultChan
//}
