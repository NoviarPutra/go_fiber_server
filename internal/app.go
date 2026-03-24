package internal

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/handlers/errors"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/routes"
)

func Bootstrap(db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Office Core API v1.0",
		ErrorHandler: errors.GlobalErrorHandler,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		BodyLimit:    4 * 1024 * 1024, // 4MB
		Prefork:      false,
	})

	// Urutan middleware sangat penting!
	app.Use(middlewares.RecoveryMiddleware())  // 1. Tangkap panic
	app.Use(middlewares.LoggerMiddleware())    // 2. Log request
	app.Use(middlewares.CorsMiddleware())      // 3. CORS
	app.Use(middlewares.RateLimitMiddleware()) // 4. Rate limit
	app.Use(middlewares.DBMiddleware(db))      // 5. Inject DB

	routes.SetupRoutes(app)

	return app
}
