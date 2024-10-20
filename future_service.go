package quasar

type Future struct {
	resultChan     <-chan []interface{}
	result         []interface{}
	executionError error
}

func NewFuture(resultChannel <-chan []interface{}) *Future {
	return &Future{
		resultChan: resultChannel,
	}
}

func (f *Future) GetResult() []interface{} {
	f.result = <-f.resultChan
	return f.result
}
