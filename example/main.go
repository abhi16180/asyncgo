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
	workerPool := executorService.NewFixedWorkerPool(&quasar.Options{
		WorkerCount: 10,
		BufferSize:  20,
	})
	defer workerPool.ShutdownGracefully()
	for i := 0; i < 10; i++ {
		f, err := workerPool.Submit(testFunction)
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

func testFunction() (myStruct, int) {
	time.Sleep(2 * time.Second)
	return myStruct{
		Value: 1,
	}, 10
}

//type Executor struct {
//}
//
//func NewExecutorService() *Executor {
//	return &Executor{}
//}
//
//func (e *Executor) Submit(fn interface{}, args ...interface{}) (*Future, error) {
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
//func (e *Executor) run(ch chan<- interface{}, fn interface{}, args ...interface{}) {
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
