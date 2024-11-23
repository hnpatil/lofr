package lofr

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gofr.dev/pkg/gofr/http"
	"reflect"
)

func validateRequest(request reflect.Value) error {
	err := validator.New().Struct(request.Interface())
	if err != nil {
		var vErr validator.ValidationErrors
		ok := errors.As(err, &vErr)
		if ok {
			return getValidationError(vErr)
		}

		return err
	}

	return nil
}

func getValidationError(errs validator.ValidationErrors) error {
	var params []string
	hasMissingParams := false

	for _, fr := range errs {
		if fr.Tag() == "required" {
			hasMissingParams = true
		}

		params = append(params, fr.Field())
	}

	if hasMissingParams {
		return http.ErrorMissingParam{Params: params}
	}

	return http.ErrorInvalidParam{Params: params}
}
