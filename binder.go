package lofr

import (
	"encoding/json"
	"errors"
	"fmt"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
	"reflect"
	"strconv"
)

type Tag string

const (
	QueryTag  Tag = "query"
	PathTag   Tag = "path"
	HeaderTag Tag = "header"
)

func bindInput(ctx *gofr.Context, ip reflect.Value) error {
	input := ip.Interface()

	err := bindJSON(ctx, input)
	if err != nil {
		return err
	}

	err = bindPathVariables(ctx, ip)
	if err != nil {
		return err
	}

	err = bindQueryParams(ctx, ip)
	if err != nil {
		return err
	}

	err = bindHeaders(ctx, ip)
	if err != nil {
		return err
	}

	return nil
}

func bindJSON(ctx *gofr.Context, input interface{}) error {
	err := ctx.Bind(input)
	if err == nil {
		return nil
	}

	var iue *json.UnmarshalTypeError

	if errors.As(err, &iue) {
		return http.ErrorInvalidParam{Params: []string{iue.Field}}
	}

	return err
}

func bindPathVariables(ctx *gofr.Context, ip reflect.Value) error {
	return bind(ip, PathTag, func(field string) (string, error) {
		return ctx.PathParam(field), nil
	})
}

func bindQueryParams(ctx *gofr.Context, ip reflect.Value) error {
	return bind(ip, QueryTag, func(field string) (string, error) {
		return ctx.Param(field), nil
	})
}

// bindHeaders won't work at the moment since gofr does not make headers accessible
func bindHeaders(ctx *gofr.Context, ip reflect.Value) error {
	return bind(ip, HeaderTag, func(field string) (string, error) {
		return "", nil
	})
}

func bind(v reflect.Value, tag Tag, extractor func(field string) (string, error)) error {
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		field := v.Field(i)

		if fieldType.Anonymous {
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
			} else {
				if field.CanAddr() {
					field = field.Addr()
				}
			}

			err := bind(field, tag, extractor)
			if err != nil {
				return err
			}

			continue
		}

		tagValue := fieldType.Tag.Get(string(tag))
		if tagValue == "" {
			continue
		}

		value, err := extractor(tagValue)
		if err != nil {
			return http.ErrorInvalidParam{Params: []string{tagValue}}
		}

		if value == "" {
			continue
		}

		err = bindValue(value, field)
		if err != nil {
			return http.ErrorInvalidParam{Params: []string{tagValue}}
		}
	}

	return nil
}

func bindValue(value string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(i)
	default:
		return fmt.Errorf("unsupported parameter type: %v", field.Kind())
	}
	return nil
}