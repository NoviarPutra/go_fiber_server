package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper struct untuk memicu validation error
type TestRequest struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=4,max=10,alphanum"`
	Password string `validate:"required,min=8"`
}

func TestFormatValidationError(t *testing.T) {
	// Tambahkan argumen (t *testing.T) di setiap callback t.Run

	t.Run("Required_Field_Error", func(t *testing.T) {
		req := TestRequest{Email: "", Username: "user", Password: "password123"}
		err := validate.Struct(req)

		msg := format_validation_error(err)
		assert.Equal(t, "Email tidak boleh kosong", msg)
	})

	t.Run("Invalid_Email_Format", func(t *testing.T) {
		req := TestRequest{Email: "bukan-email", Username: "user", Password: "password123"}
		err := validate.Struct(req)

		msg := format_validation_error(err)
		assert.Equal(t, "Format email tidak valid", msg)
	})

	t.Run("Min_Length_Constraint", func(t *testing.T) {
		req := TestRequest{Email: "test@officecore.id", Username: "abc", Password: "password123"}
		err := validate.Struct(req)

		msg := format_validation_error(err)
		assert.Equal(t, "Username minimal 4 karakter", msg)
	})

	t.Run("Max_Length_Constraint", func(t *testing.T) {
		req := TestRequest{Email: "test@officecore.id", Username: "usernameterlalupanjang", Password: "password123"}
		err := validate.Struct(req)

		msg := format_validation_error(err)
		assert.Equal(t, "Username maksimal 10 karakter", msg)
	})

	t.Run("Alphanum_Constraint_Violation", func(t *testing.T) {
		req := TestRequest{Email: "test@officecore.id", Username: "user-123", Password: "password123"}
		err := validate.Struct(req)

		msg := format_validation_error(err)
		assert.Equal(t, "Username hanya boleh huruf dan angka", msg)
	})

	t.Run("Non_Validator_Error_Fallback", func(t *testing.T) {
		err := errors.New("random db error")

		msg := format_validation_error(err)
		assert.Equal(t, "Input tidak valid", msg)
	})

	t.Run("Nil_Error_Safety", func(t *testing.T) {
		msg := format_validation_error(nil)
		assert.Equal(t, "Input tidak valid", msg)
	})
}
