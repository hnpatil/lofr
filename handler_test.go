package lofr

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"testing"
)

func Test_validateHandler(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		err   error
	}{
		{
			name:  "case: input = (*gofr.Context), op = (error)",
			input: func(ctx *gofr.Context) error { return nil },
			err:   nil,
		},
		{
			name:  "case: input = (*gofr.Context, [object ptr]), op = (error)",
			input: func(ctx *gofr.Context, input *struct{}) error { return nil },
			err:   nil,
		},
		{
			name:  "case: input = (*gofr.Context, [object ptr]), op = ([object ptr], error)",
			input: func(ctx *gofr.Context, input *struct{}) (*struct{}, error) { return nil, nil },
			err:   nil,
		},
		{
			name:  "case: input = (*gofr.Context, [object]), op = ([object ptr], error)",
			input: func(ctx *gofr.Context, input struct{}) (*struct{}, error) { return nil, nil },
			err:   errors.New(""),
		},
		{
			name:  "case: input = (*gofr.Context, [object ptr]), op = ([object], error)",
			input: func(ctx *gofr.Context, input *struct{}) (struct{}, error) { return struct{}{}, nil },
			err:   errors.New(""),
		},
		{
			name:  "case: input = (*gofr.Context, [object ptr], [object ptr]), op = ([object ptr], error)",
			input: func(ctx *gofr.Context, input1 *struct{}, input2 *struct{}) (*struct{}, error) { return nil, nil },
			err:   errors.New(""),
		},
		{
			name:  "case: input = (*gofr.Context, [object ptr]), op = ([object ptr], [object ptr], error)",
			input: func(ctx *gofr.Context, input struct{}) (*struct{}, *struct{}, error) { return nil, nil, nil },
			err:   errors.New(""),
		},
		{
			name:  "case: input = ([object ptr]), op = (error)",
			input: func(input struct{}) error { return nil },
			err:   errors.New(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok && test.err == nil {
						t.Errorf("was not expecting error, got: %v", err)
					}
				}
			}()

			Handler(test.input)
		})
	}
}
