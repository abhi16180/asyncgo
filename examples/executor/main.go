package main

import (
	"fmt"
	"github.com/abhi16180/asyncgo"
	"log"
	"time"
)

func main() {
	now := time.Now()
	executor := asyncgo.NewExecutor()
	// submit any function signature
	// first param is  function, subsequent params are arguments
	future1 := executor.Submit(func(arg1, arg2 int) int {
		// mocking delay
		time.Sleep(1 * time.Second)
		return arg1 + arg2
	}, 10, 20)

	future2 := executor.Submit(func(arg1 int) int {
		// mocking delay
		time.Sleep(1 * time.Second)
		return arg1 * 10
	}, 20)

	// getting back the results
	result1, err := future1.Get()
	if err != nil {
		log.Println(err)
	}
	result2, err := future2.Get()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result1, result2)
	fmt.Println("time taken %v", time.Since(now))

}
