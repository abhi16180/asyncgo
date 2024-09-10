package wp

// TODO add more funcs
type Future struct {
	resultChan <-chan []interface{}
	result     []interface{}
}

func NewFuture(resultChannel <-chan []interface{}) *Future {
	return &Future{
		resultChan: resultChannel,
	}
}

func (f *Future) GetResult() interface{} {
	return <-f.resultChan
}
