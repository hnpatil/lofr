package lofr

import (
	"gofr.dev/pkg/gofr/http"
	"reflect"
)

func setDefaults(v reflect.Value) error {
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		field := v.Field(i)

		if fieldType.Anonymous {
			err := setAnonymousDefaults(field)
			if err != nil {
				return err
			}

			continue
		}

		tagValue := fieldType.Tag.Get(string(DefaultTag))
		if tagValue == "" {
			continue
		}

		if !field.IsZero() {
			continue
		}

		err := bindValue(tagValue, field)
		if err != nil {
			return http.ErrorInvalidParam{Params: []string{tagValue}}
		}
	}

	return nil
}

func setAnonymousDefaults(field reflect.Value) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else {
		if field.CanAddr() {
			field = field.Addr()
		}
	}

	err := setDefaults(field)
	if err != nil {
		return err
	}

	return nil
}
