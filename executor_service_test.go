package wp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
			err:     fmt.Errorf("function must be a function"),
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
			e := &ExecutorServiceImpl{}
			got, _ := e.Submit(tt.args.function, tt.args.args...)
			result, err := got.GetResult()
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
			e := &ExecutorServiceImpl{}
			wp := e.NewFixedWorkerPool(tt.args.options)
			assert.NotNil(t, wp, "NewFixedWorkerPool(%v)", tt.args.options)
			wp.Terminate()
		})
	}
}
