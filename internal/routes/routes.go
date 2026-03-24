package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/internal/handlers/auth"
	"github.com/yourusername/go_server/internal/handlers/base"
	"github.com/yourusername/go_server/internal/handlers/health"
	"github.com/yourusername/go_server/internal/handlers/users"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/utils"
)

var not_found = func(c *fiber.Ctx) error {
	return utils.NotFound(c, "Route tidak ditemukan")
}

func SetupRoutes(app *fiber.App) {
	// ─── Public ───────────────────────────────────────────────────────────────
	app.Get("/", base.InitHandler)
	app.Get("/health", health.HealthHandler)

	// ─── API v1 ───────────────────────────────────────────────────────────────
	api := app.Group("/api/v1")

	// Auth
	auth_group := api.Group("/auth")
	auth_group.Post("/register", auth.Register)
	auth_group.Post("/login", auth.Login)

	// Users (private)
	users_group := api.Group("/users")
	users_group.Get("/", middlewares.Protected(), middlewares.Pagination, users.GetAll)

	// 404
	app.Use(not_found)
}
