package helper

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}

func ValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, fieldError := range validationErrors {
		field := fieldError.Field()

		switch fieldError.Tag() {
		case "required":
			errors[field] = field + " is required"

		case "email":
			errors[field] = field + " must be a valid email"

		case "min":
			errors[field] = field + " does not meet the minimum length"

		case "max":
			errors[field] = field + " exceeds the maximum length"

		default:
			errors[field] = field + " is invalid"
		}
	}

	return errors
}
