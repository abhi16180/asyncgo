package asyncgo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecutorServiceImpl_Submit(t *testing.T) {
	type args struct {
		function interface{}
		args     []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
		err     error
	}{
		{
			name: "success",
			args: args{
				function: func() (interface{}, error) {
					return 10, nil
				},
			},
			want: []interface{}{
				10, nil,
			},
			wantErr: false,
		},
		{
			name: "fails due to invalid function",
			args: args{
				function: "wrongParam",
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("parameter 'function' must be a function"),
		},
		{
			name: "fails due to invalid args",
			args: args{
				function: func(a int, b int) (interface{}, error) {
					return a + b, nil
				},
				args: []interface{}{},
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("function must have %d parameters", 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExecutorService{}
			got := e.Submit(tt.args.function, tt.args.args...)
			result, err := got.Get()
			if tt.wantErr {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestExecutorServiceImpl_NewFixedWorkerPool(t *testing.T) {
	type args struct {
		options *Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				options: &Options{
					WorkerCount: 2,
					BufferSize:  10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExecutorService{}
			wp := e.NewFixedWorkerPool(context.Background(), tt.args.options)
			assert.NotNil(t, wp, "NewFixedWorkerPool(%v)", tt.args.options)
			wp.Shutdown()
		})
	}
}

func TestWorkerPool1(t *testing.T) {
	executorService := NewExecutor()
	workerPool := executorService.NewFixedWorkerPool(context.TODO(), &Options{
		WorkerCount: 100,
		BufferSize:  100,
	})
	multiply := func(a, b int) int {
		time.Sleep(time.Second)
		return a * b
	}
	futures := make([]*Future, 0)
	expectedSlice := make([]int, 0)
	for i := 0; i < 100; i++ {
		expected := i * (i + 1)
		f, err := workerPool.Submit(multiply, i, i+1)
		futures = append(futures, f)
		expectedSlice = append(expectedSlice, expected)
		if err != nil {
			return
		}
	}
	for i, future := range futures {
		result, err := future.Get()
		assert.Nil(t, err)
		assert.Equal(t, result[0].(int), expectedSlice[i])
	}
	workerPool.Shutdown()
}

func TestWorkerPoolEnsureGracefulShutdown(t *testing.T) {
	executorService := NewExecutor()
	workerPool := executorService.NewFixedWorkerPool(context.TODO(), &Options{
		WorkerCount: 100,
		BufferSize:  100,
	})
	multiply := func(a, b int) int {
		time.Sleep(time.Second)
		return a * b
	}
	futures := make([]*Future, 0)
	expectedSlice := make([]int, 0)
	for i := 0; i < 100; i++ {
		expected := i * (i + 1)
		f, err := workerPool.Submit(multiply, i, i+1)
		futures = append(futures, f)
		expectedSlice = append(expectedSlice, expected)
		if err != nil {
			return
		}
	}
	// even though shutdown is called even before retrieving results, it should not cancel the existing tasks
	// all the submitted tasks' execution should be guaranteed.
	workerPool.Shutdown()
	for i, future := range futures {
		result, err := future.Get()
		assert.Nil(t, err)
		assert.Equal(t, result[0].(int), expectedSlice[i])
	}
}

func TestWorkerPoolSubmitTaskAfterShutdown(t *testing.T) {
	executorService := NewExecutor()
	workerPool := executorService.NewFixedWorkerPool(context.TODO(), nil)
	multiply := func(a, b int) int {
		time.Sleep(time.Second)
		return a * b
	}
	futures := make([]*Future, 0)
	expectedSlice := make([]int, 0)
	for i := 0; i < 10; i++ {
		expected := i * (i + 1)
		f, err := workerPool.Submit(multiply, i, i+1)
		futures = append(futures, f)
		expectedSlice = append(expectedSlice, expected)
		if err != nil {
			return
		}
	}
	workerPool.Shutdown()
	for i, future := range futures {
		result, err := future.Get()
		assert.Nil(t, err)
		assert.Equal(t, result[0].(int), expectedSlice[i])
	}
	f, err := workerPool.Submit(func() int {
		time.Sleep(time.Second)
		return 10
	})
	assert.Nil(t, f)
	assert.Equal(t, fmt.Errorf("cannot add new task after closing worker pool"), err)

}
