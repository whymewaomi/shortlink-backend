package validator

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func Valide(s interface{}) error {
	return Validate.Struct(s)
}