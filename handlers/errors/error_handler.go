package errors

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/yourusername/go_server/utils"
)

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	// 1. Tampilkan Stack Trace jika error berasal dari pkg/errors
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err != nil {
		log.Printf("[ERROR] %v", err)

		// Mengecek apakah error memiliki stack trace
		if e, ok := err.(stackTracer); ok {
			log.Printf("Stack Trace: %+v", e.StackTrace())
		}
	}

	// 2. Logic tetap sama
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return utils.SendResponse[any](c, code, false, message, nil, nil)
}
