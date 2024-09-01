package lofr

import (
	"fmt"
	"reflect"
	"runtime"

	"gofr.dev/pkg/gofr"
)

// Handler returns a gofr.Handler that wraps the handler passed as parameter.
// handler may use the following signature:
//
//	func(*gofr.Context, [input object ptr]) ([output object], error)
//
// Input and output objects are both optional.
// As such, the minimal accepted signature is:
//
//	func(*gofr.Context) error
//
// The wrapper will bind input object fields with the following tags
// - json :- fields being passed in the body of the request
// - path :- fields being passed as path variables in the request
// - query :- fields being passed as query param in the request
//
// Handler will panic if the handler is incompatible.
func Handler(handler interface{}) gofr.Handler {
	hv := reflect.ValueOf(handler)

	ip, op, err := validateHandler(hv)
	if err != nil {
		panic(err)
	}

	return func(c *gofr.Context) (interface{}, error) {
		args := []reflect.Value{reflect.ValueOf(c)}

		if ip != nil {
			input := reflect.New(ip)
			err = bindInput(c, input)
			if err != nil {
				return nil, err
			}

			args = append(args, input)
		}

		var (
			response interface{}
			hError   interface{}
		)

		ret := hv.Call(args)

		if op != nil {
			response = ret[0].Interface()
			hError = ret[1].Interface()
		} else {
			hError = ret[0].Interface()
		}

		if hError == nil {
			return response, nil
		}

		return nil, hError.(error)
	}
}

func validateHandler(hv reflect.Value) (in, op reflect.Type, err error) {
	if hk := hv.Kind(); hk != reflect.Func {
		return nil, nil, fmt.Errorf("handler parameters must be a function, got %s", hk.String())
	}

	op, err = validateOutput(hv)
	if err != nil {
		return nil, nil, err
	}

	in, err = validateInput(hv)
	if err != nil {
		return nil, nil, err
	}

	return in, op, nil
}

func validateOutput(hv reflect.Value) (reflect.Type, error) {
	ht := hv.Type()
	n := ht.NumOut()
	fName := fmt.Sprintf("%s", runtime.FuncForPC(hv.Pointer()).Name())

	if n != 1 && n != 2 {
		return nil, fmt.Errorf("handler %s parameters must return one or two return values, got %d", fName, n)
	}

	if !ht.Out(n - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, fmt.Errorf("handler %s does not return an error, got %s", fName, ht.Out(n-1).String())
	}

	if n == 2 {
		t := ht.Out(0)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return t, nil
	}

	return nil, nil
}

func validateInput(hv reflect.Value) (reflect.Type, error) {
	ht := hv.Type()
	n := ht.NumIn()
	fName := fmt.Sprintf("%s", runtime.FuncForPC(hv.Pointer()).Name())

	if n != 1 && n != 2 {
		return nil, fmt.Errorf("handler %s parameters must accept one or two parameters, got %d", fName, n)
	}

	if !ht.In(0).ConvertibleTo(reflect.TypeOf(&gofr.Context{})) {
		return nil, fmt.Errorf("handler %s does not accept *gofr.Context as first parameter, got %s", fName, ht.In(0).String())
	}

	if n == 2 {
		if ht.In(1).Kind() != reflect.Ptr || ht.In(1).Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("handler %s does not accept pointer to struct as second parameter, got %s", fName, ht.In(1).String())
		} else {
			return ht.In(1).Elem(), nil
		}
	}

	return nil, nil
}
