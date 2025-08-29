package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GetDetailedError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		firstError := validationErrors[0]

		switch firstError.Tag() {
		case "required":
			return fmt.Errorf("field %s required", firstError.Field())
		case "min":
			return fmt.Errorf("field %s should bigger then %s", firstError.Field(), firstError.Param())
		case "gt":
			return fmt.Errorf("field %s should greater then %s", firstError.Field(), firstError.Param())
		default:
			return fmt.Errorf("field %s is incorrect", firstError.Field())
		}
	}
	return fmt.Errorf("validation error")
}

func Validate(s interface{}) error {
	return validator.New().Struct(s)
}
