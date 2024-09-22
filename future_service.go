package wp

// TODO add more funcs

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

func (f *Future) GetResult() ([]interface{}, error) {
	f.result = <-f.resultChan
	if len(f.result) == 1 {
		switch v := f.result[0].(type) {
		case error:
			f.executionError = v
		}
	}
	return f.result, f.executionError
}
