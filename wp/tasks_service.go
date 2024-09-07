package wp

import "reflect"

type Task interface {
	Execute(function interface{}, args []interface{}, resultChannel chan<- interface{}) error
}

type TaskImpl struct {
}

func NewTask() Task {
	return &TaskImpl{}
}

func (t *TaskImpl) Execute(function interface{}, args []interface{}, resultChannel chan<- interface{}) error {
	val := reflect.ValueOf(function)
	argSlice := make([]reflect.Value, len(args))
	for i, arg := range args {
		argSlice[i] = reflect.ValueOf(arg)
	}
	if len(argSlice) > 0 {
		result := val.Call(argSlice)
		resultChannel <- result
		return nil
	}
	resultChannel <- val.Call(argSlice)
	return nil
}
