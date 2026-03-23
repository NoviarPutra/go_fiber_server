package auth

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
		case "email":
			return "Format email tidak valid"
		case "min":
			return field.Field() + " minimal " + field.Param() + " karakter"
		case "max":
			return field.Field() + " maksimal " + field.Param() + " karakter"
		case "alphanum":
			return field.Field() + " hanya boleh huruf dan angka"
		}
	}
	return "Input tidak valid"
}
