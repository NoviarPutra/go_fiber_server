package middlewares

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/yourusername/go_server/internal/utils"
)

func RateLimitMiddleware() fiber.Handler {
	max_requests := 100
	if os.Getenv("APP_ENV") == "development" {
		max_requests = 1000
	}

	return limiter.New(limiter.Config{
		Max:        max_requests,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests,
				"Terlalu banyak request. Coba lagi dalam 1 menit.",
			)
		},
	})
}
