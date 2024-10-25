package commons

import "time"

type Options struct {
	// WorkerCount - number of goroutines
	WorkerCount int64
	// BufferSize is size of buffered channel
	BufferSize int64
	// IdleSleepDuration is needed to specify sleep duration if not new tasks is added to queue.
	// It is required to prevent unnecessary wasting of CPU cycles.
	// Default value will be set as 10ms
	IdleSleepDuration time.Duration
}
