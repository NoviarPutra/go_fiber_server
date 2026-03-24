package middlewares

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RecoveryMiddleware() fiber.Handler {
	is_dev := os.Getenv("APP_ENV") == "development"

	return recover.New(recover.Config{
		EnableStackTrace: is_dev,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			if is_dev {
				log.Printf("🔥 PANIC RECOVERED: %v", e)
			}
		},
	})
}
