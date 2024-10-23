package asyncgo

import (
	"fmt"
	"github.com/abhi16180/asyncgo/utils"
	"reflect"
)

//go:generate mockery --name=Task --output=./mocks --outpkg=mocks
type Task interface {
	// Execute gets the function signature using reflection. Calls the function
	Execute() error
}

type TaskService struct {
	resultChannel chan<- []interface{}
	errChan       chan<- error
	function      interface{}
	args          []interface{}
}

func NewTask(resultChan chan<- []interface{}, errChannel chan<- error, function interface{}, args []interface{}) Task {
	return &TaskService{
		resultChannel: resultChan,
		errChan:       errChannel,
		function:      function,
		args:          args,
	}
}

func (t *TaskService) Execute() error {
	val := reflect.ValueOf(t.function)
	kind := val.Kind()
	if kind != reflect.Func {
		t.resultChannel <- nil
		t.errChan <- fmt.Errorf("parameter 'function' must be a function")
		return fmt.Errorf("parameter 'function' must be a function")
	}
	numIn := val.Type().NumIn()
	if numIn != len(t.args) {
		t.resultChannel <- nil
		t.errChan <- fmt.Errorf("function must have %d parameters", numIn)
		return fmt.Errorf("function must have %d parameters", numIn)
	}
	argSlice := make([]reflect.Value, len(t.args))
	for i, arg := range t.args {
		argSlice[i] = reflect.ValueOf(arg)
	}
	var result []reflect.Value
	if len(argSlice) > 0 {
		result = val.Call(argSlice)
	} else {
		result = val.Call([]reflect.Value{})
	}
	t.resultChannel <- utils.GetResultInterface(result)
	t.errChan <- nil
	return nil
}
