package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/handlers/base"
	"github.com/yourusername/go_server/handlers/health"
	"github.com/yourusername/go_server/handlers/users"
	"github.com/yourusername/go_server/middlewares"
	"github.com/yourusername/go_server/utils"
)

var not_found = func(c *fiber.Ctx) error {
	return utils.NotFound(c, "Route tidak ditemukan")
}

func SetupRoutes(app *fiber.App) {
	// Public Routes
	app.Get("/", base.InitHandler)
	app.Get("/health", health.HealthHandler)

	// API Group
	api := app.Group("/api/v1")

	// Users Group (Private)
	users_group := api.Group("/users")
	users_group.Get("/", middlewares.Protected(), middlewares.Pagination, users.UsersHandler)

	// 404
	app.Use(not_found)
}
