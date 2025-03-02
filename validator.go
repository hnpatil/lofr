package lofr

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gofr.dev/pkg/gofr/http"
	"reflect"
)

var skipValidation bool

// SkipValidation sets skipValidation. When set as true, request validation is skipped.
func SkipValidation(value bool) {
	skipValidation = value
}

func validateRequest(request reflect.Value) error {
	if skipValidation {
		return nil
	}

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
