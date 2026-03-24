package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsMiddleware() fiber.Handler {
	allowed_origins := os.Getenv("ALLOWED_ORIGINS")

	// Development: izinkan semua origin (Postman, localhost, dll)
	if os.Getenv("APP_ENV") == "development" || allowed_origins == "" {
		allowed_origins = "*"
	}

	return cors.New(cors.Config{
		AllowOrigins: allowed_origins,
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type,Authorization",
		MaxAge:       86400,
	})
}
