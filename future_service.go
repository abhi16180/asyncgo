package asyncgo

type Future struct {
	resultChan     <-chan []interface{}
	errChan        <-chan error
	result         []interface{}
	err            error
	executionError error
}

func NewFuture(resultChannel <-chan []interface{}, errChan chan error) *Future {
	return &Future{
		resultChan: resultChannel,
		errChan:    errChan,
	}
}

func (f *Future) Get() ([]interface{}, error) {
	f.result = <-f.resultChan
	f.err = <-f.errChan
	return f.result, f.err
}

func (f *Future) Wait() error {
	_, err := f.Get()
	return err
}

// wait all futures
