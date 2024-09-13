package wp

import (
	"reflect"
	"wp/utils"
)

//go:generate mockery --name=Task --output=./mocks --outpkg=mocks
type Task interface {
	// Execute gets the function signature using reflection. Calls the function
	Execute() error
}

type TaskImpl struct {
	resultChannel chan<- []interface{}
	function      interface{}
	args          []interface{}
}

func NewTask(resultChan chan<- []interface{}, function interface{}, args []interface{}) Task {
	return &TaskImpl{
		resultChannel: resultChan,
		function:      function,
		args:          args,
	}
}

func (t *TaskImpl) Execute() error {
	val := reflect.ValueOf(t.function)
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
	return nil
}
