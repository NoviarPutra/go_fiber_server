package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/handlers/base"
	"github.com/yourusername/go_server/handlers/health"
	"github.com/yourusername/go_server/handlers/users"
	"github.com/yourusername/go_server/middlewares"
)

func SetupRoutes(app *fiber.App) {
	// Route setup logic here
	app.Get("/", base.InitHandler)
	app.Get("/health", health.HealthHandler)

	api := app.Group("/api/v1")

	usersRoute := api.Group("/users")
	usersRoute.Use(middlewares.Protected())
	usersRoute.Get("/", middlewares.Pagination, users.UsersHandler)
}
