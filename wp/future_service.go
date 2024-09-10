package wp

type Future struct {
	result <-chan []interface{}
}

func NewFuture(resultChannel <-chan []interface{}) *Future {
	return &Future{
		result: resultChannel,
	}
}

func (f *Future) Result() interface{} {
	return <-f.result
}
