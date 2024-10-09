package quasar

import (
	"fmt"
	"github.com/abhi16180/quasar/utils"
	"reflect"
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
	kind := val.Kind()
	if kind != reflect.Func {
		t.resultChannel <- []interface{}{fmt.Errorf("function must be a function")}
		return fmt.Errorf("function must be a function")
	}
	numIn := val.Type().NumIn()
	if numIn != len(t.args) {
		t.resultChannel <- []interface{}{fmt.Errorf("function must have %d parameters", numIn)}
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
	return nil
}
