package lofr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
	"net/http/httptest"
	"net/url"
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

func Test_QueryParams(t *testing.T) {
	tests := []struct {
		desc    string
		queries map[string]interface{}
		res     *basicRequest
		err     error
	}{
		{
			desc: "case: valid query values",
			queries: map[string]interface{}{
				"name":     "User Name",
				"age":      27,
				"score":    4.78,
				"balance":  -100,
				"isActive": true,
			},
			res: &basicRequest{
				Name:     "User Name",
				Age:      27,
				Score:    4.78,
				Balance:  -100,
				IsActive: true,
			},
		},
		{
			desc: "case: invalid query value for uint",
			queries: map[string]interface{}{
				"name": "User Name",
				"age":  "User Name",
			},
			err: http.ErrorInvalidParam{Params: []string{"age"}},
		},
		{
			desc: "case: invalid query value for bool",
			queries: map[string]interface{}{
				"name":     "User Name",
				"isActive": "User Name",
			},
			err: http.ErrorInvalidParam{Params: []string{"isActive"}},
		},
		{
			desc: "case: invalid query value for float",
			queries: map[string]interface{}{
				"name":  "User Name",
				"score": "User Name",
			},
			err: http.ErrorInvalidParam{Params: []string{"score"}},
		},
		{
			desc: "case: invalid query value for int",
			queries: map[string]interface{}{
				"name":    "User Name",
				"balance": "User Name",
			},
			err: http.ErrorInvalidParam{Params: []string{"balance"}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			queryParams := url.Values{}

			for k, v := range test.queries {
				queryParams.Add(k, fmt.Sprintf("%v", v))
			}

			u := &url.URL{
				Path: "/basic",
			}

			u.RawQuery = queryParams.Encode()

			req := httptest.NewRequest("GET", u.String(), nil)

			ctx := &gofr.Context{
				Context: context.TODO(),
				Request: http.NewRequest(req),
			}

			res, err := Handler(basicHandler)(ctx)
			assert.Equal(t, test.err, err)
			if err != nil {
				assert.Nil(t, test.res)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func Test_RequestBody(t *testing.T) {
	tests := []struct {
		desc string
		body []byte
		res  *basicRequest
		err  error
	}{
		{
			desc: "case: valid request body",
			body: []byte(`{"firstName":"First Name","lastName":"Last Name"}`),
			res: &basicRequest{
				Person: Person{
					FirstName: "First Name",
					LastName:  "Last Name",
				},
			},
		},
		{
			desc: "case: invalid filed type in body",
			body: []byte(`{"id":"First Name","lastName":"Last Name"}`),
			err:  http.ErrorInvalidParam{Params: []string{"id"}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/basic", bytes.NewBuffer(test.body))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{
				Context: context.TODO(),
				Request: http.NewRequest(req),
			}

			res, err := Handler(basicHandler)(ctx)
			assert.Equal(t, test.err, err)
			if err != nil {
				assert.Nil(t, test.res)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ID        int    `json:"id"`
}

type basicRequest struct {
	Person
	Name     string  `query:"name"`
	Age      uint    `query:"age"`
	Score    float64 `query:"score"`
	Balance  int     `query:"balance"`
	IsActive bool    `query:"isActive"`
}

func basicHandler(ctx *gofr.Context, req *basicRequest) (*basicRequest, error) {
	return req, nil
}

func Test_QueryParamsWithDefaultValues(t *testing.T) {
	tests := []struct {
		desc    string
		queries map[string]interface{}
		res     *queryWithDefaults
		err     error
	}{
		{
			desc:    "case: no queries passed",
			queries: map[string]interface{}{},
			res: &queryWithDefaults{
				Name:    "person",
				Age:     21,
				Score:   8.5,
				Balance: 100,
			},
		},
		{
			desc: "case: some queries passed",
			queries: map[string]interface{}{
				"name": "User Name",
				"age":  25,
			},
			res: &queryWithDefaults{
				Name:    "User Name",
				Age:     25,
				Score:   8.5,
				Balance: 100,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			queryParams := url.Values{}

			for k, v := range test.queries {
				queryParams.Add(k, fmt.Sprintf("%v", v))
			}

			u := &url.URL{
				Path: "/basic",
			}

			u.RawQuery = queryParams.Encode()

			req := httptest.NewRequest("GET", u.String(), nil)

			ctx := &gofr.Context{
				Context: context.TODO(),
				Request: http.NewRequest(req),
			}

			res, err := Handler(defaultQueriesHandler)(ctx)
			assert.Equal(t, test.err, err)
			if err != nil {
				assert.Nil(t, test.res)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

type queryWithDefaults struct {
	Name    string  `query:"name" default:"person"`
	Age     uint    `query:"age" default:"21"`
	Score   float64 `query:"score" default:"8.5"`
	Balance int     `query:"balance" default:"100"`
}

func defaultQueriesHandler(ctx *gofr.Context, req *queryWithDefaults) (*queryWithDefaults, error) {
	return req, nil
}

func Test_PostWithValidation(t *testing.T) {
	tests := []struct {
		desc string
		body []byte
		err  error
	}{
		{
			desc: "case: missing email",
			body: []byte(`{"name":"First Name","age":21}`),
			err:  http.ErrorMissingParam{Params: []string{"Email"}},
		},
		{
			desc: "case: invalid email",
			body: []byte(`{"name":"First Name","age":21, "email": "1234567890"}`),
			err:  http.ErrorInvalidParam{Params: []string{"Email"}},
		},
		{
			desc: "case: invalid request",
			body: []byte(`{"name":"First Name","age":21, "email": "lite@gofr.com"}`),
		},
		{
			desc: "case: invalid age",
			body: []byte(`{"name":"First Name","age":18, "email": "lite@gofr.com"}`),
			err:  http.ErrorInvalidParam{Params: []string{"Age"}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/basic", bytes.NewBuffer(test.body))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{
				Context: context.TODO(),
				Request: http.NewRequest(req),
			}

			_, err := Handler(defaultPostHandler)(ctx)
			assert.Equal(t, test.err, err)
		})
	}
}

type postRequest struct {
	Name  string `json:"name"`
	Age   uint   `json:"age" validate:"gte=21"`
	Email string `json:"email" validate:"required,email"`
}

func defaultPostHandler(ctx *gofr.Context, req *postRequest) error {
	return nil
}
