package validation

import (
	"github.com/go-playground/validator/v10"
)

func ValidateData(input any) map[string]string {
	validate := validator.New()
	validationError := validate.Struct(input)

	allErrors := make(map[string]string)

	if validationError != nil {
		for _, err := range validationError.(validator.ValidationErrors) {
			allErrors[err.Field()] = err.Error()
		}

		return allErrors
	}

	return nil
}
