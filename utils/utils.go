package utils

import "reflect"

func GetResultInterface(result []reflect.Value) []interface{} {
	resultInterface := make([]interface{}, 0)
	for _, item := range result {
		resultInterface = append(resultInterface, item.Interface())
	}
	return resultInterface
}
