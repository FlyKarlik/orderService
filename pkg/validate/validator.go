package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(r any) error {
	return validate.Struct(r)
}
