package errors

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	pkg_errors "github.com/pkg/errors"
	"github.com/yourusername/go_server/utils"
)

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	is_dev := os.Getenv("APP_ENV") == "development"

	// 1. Log error selalu
	log.Printf("[ERROR] %v", err)

	// 2. Stack trace hanya di development
	if is_dev {
		type stack_tracer interface {
			StackTrace() pkg_errors.StackTrace
		}
		if e, ok := err.(stack_tracer); ok {
			log.Printf("[STACK TRACE]\n%+v", e.StackTrace())
		}
	}

	// 3. Default: 500
	code := fiber.StatusInternalServerError
	message := "Terjadi kesalahan pada server"

	// 4. Override jika Fiber error (404, 405, dll)
	var fiber_err *fiber.Error
	if errors.As(err, &fiber_err) {
		code = fiber_err.Code
		message = fiber_err.Message
	}

	return utils.ErrorResponse(c, code, message)
}
