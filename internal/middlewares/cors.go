package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsMiddleware() fiber.Handler {
	allowed_origins := os.Getenv("ALLOWED_ORIGINS")

	// Development: Bullet-proof CORS (tidak boleh "*" jika AllowCredentials: true)
	if os.Getenv("APP_ENV") == "development" || allowed_origins == "" {
		allowed_origins = "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://localhost:4173"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowed_origins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
		MaxAge:           86400,
	})
}
