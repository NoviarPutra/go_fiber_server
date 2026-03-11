package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/handlers/base"
	"github.com/yourusername/go_server/handlers/health"
	"github.com/yourusername/go_server/handlers/users"
	"github.com/yourusername/go_server/middlewares"
)

func SetupRoutes(app *fiber.App) {
	// Public Routes
	app.Get("/", base.InitHandler)
	app.Get("/health", health.HealthHandler)

	// API Group
	api := app.Group("/api/v1")

	// Users Group (Private)
	usersGroup := api.Group("/users", middlewares.Protected())
	{
		usersGroup.Get("/", middlewares.Pagination, users.UsersHandler)
	}
}
