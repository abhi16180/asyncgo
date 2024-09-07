package wp

import "sync"

type Worker interface {
}

type WorkerImpl struct {
}

func NewWorker(wg *sync.WaitGroup, in <-chan interface{}, out chan<- Worker) {
	defer wg.Done()
}
