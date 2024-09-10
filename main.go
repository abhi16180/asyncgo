package main

import (
	"fmt"
	"time"
	"wp/wp"
)

func main() {
	//executorService := NewExecutorService()
	futures := make([]*wp.Future, 0)
	//
	//for i := 0; i < 1235; i++ {
	//	f, err := executorService.Submit(s)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	futures = append(futures, f)
	//}
	//
	//for _, future := range futures {
	//	fmt.Println(future.Get())
	//}
	fmt.Println("Hello World")

	executorService := wp.NewExecutorService()
	workerPool := executorService.NewFixedWorkerPool(10)
	for i := 0; i < 100; i++ {
		f, _ := workerPool.Submit(s)
		futures = append(futures, f)
	}
	for i := 0; i < 100; i++ {
		fmt.Println(i, futures[i].Result())
	}
}

func s() int {
	time.Sleep(2000 * time.Millisecond)
	return 63
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
