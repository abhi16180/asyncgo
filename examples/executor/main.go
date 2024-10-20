package main

import (
	"fmt"
	"github.com/abhi16180/asyncgo"
	"log"
	"time"
)

var mockSleepDuration = time.Second * 2

func main() {

	now := time.Now()

	// create new executor
	executor := asyncgo.NewExecutor()

	// submit any function signature
	// first param is  function, subsequent params are arguments
	future1 := executor.Submit(func(arg1, arg2 int) int {
		// mocking delay
		time.Sleep(mockSleepDuration)
		return arg1 + arg2
	}, 10, 20)

	future2 := executor.Submit(func(arg1 int) int {
		// mocking delay
		time.Sleep(mockSleepDuration)
		return arg1 * 10
	}, 20)

	// you can define a function somewhere and provide function reference with args to execute it asynchronously
	future3 := executor.Submit(someLongTask, 10)

	// NOTE - err returned by future.Get() represents error that was encountered while executing your function.
	// It does not represent the error returned by your function
	// To access error returned by your function you need to convert interface to error type from the result []interface

	result1, err := future1.Get()
	if err != nil {
		log.Println(err)
	}
	result2, err := future2.Get()
	if err != nil {
		log.Println(err)
	}

	result3, err := future3.Get()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result1, result2, result3)
	fmt.Println("time taken %v", time.Since(now))
}

func someLongTask(value int) int {
	time.Sleep(mockSleepDuration)
	return value
}
