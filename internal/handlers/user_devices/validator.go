package user_devices

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func format_validation_error(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		field := ve[0]
		switch field.Tag() {
		case "required":
			return field.Field() + " tidak boleh kosong"
		case "oneof":
			return field.Field() + " harus salah satu dari " + field.Param()
		}
	}
	return "Input tidak valid"
}
