package wp

import "reflect"

type Task interface {
	Execute() error
}

type TaskImpl struct {
	resultChannel chan<- interface{}
	function      interface{}
	args          []interface{}
}

func NewTask(resultChan chan<- interface{}, function interface{}, args []interface{}) Task {
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
	if len(argSlice) > 0 {
		result := val.Call(argSlice)
		t.resultChannel <- result
		return nil
	}
	t.resultChannel <- val.Call([]reflect.Value{})
	return nil
}
