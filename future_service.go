package asyncgo

import "fmt"

type Future struct {
	resultChan     <-chan []interface{}
	errChan        <-chan error
	result         []interface{}
	err            error
	executionError error
	isRead         bool
}

func NewFuture(resultChannel <-chan []interface{}, errChan chan error) *Future {
	return &Future{
		resultChan: resultChannel,
		errChan:    errChan,
	}
}

func (f *Future) Get() ([]interface{}, error) {
	if f.isRead {
		return f.result, f.executionError
	}
	f.result = <-f.resultChan
	f.err = <-f.errChan
	f.isRead = true
	return f.result, f.err
}

func (f *Future) Wait() error {
	if !f.isRead {
		return fmt.Errorf("asyncgo.Future: wait already read")
	}
	_, err := f.Get()
	return err
}
