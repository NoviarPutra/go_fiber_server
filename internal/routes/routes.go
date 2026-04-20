package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/internal/handlers/audit_logs"
	"github.com/yourusername/go_server/internal/handlers/auth"
	"github.com/yourusername/go_server/internal/handlers/base"
	"github.com/yourusername/go_server/internal/handlers/companies"
	"github.com/yourusername/go_server/internal/handlers/company_branches"
	"github.com/yourusername/go_server/internal/handlers/company_users"
	"github.com/yourusername/go_server/internal/handlers/divisions"
	"github.com/yourusername/go_server/internal/handlers/health"
	"github.com/yourusername/go_server/internal/handlers/user_devices"
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
	auth_group.Post("/refresh", auth.Refresh)
	auth_group.Post("/revoke", auth.Revoke)

	// Users (private)
	users_group := api.Group("/users")
	users_group.Get("/", middlewares.Protected(), middlewares.Pagination, users.GetAll)

	// Companies (private)
	companies_group := api.Group("/companies")
	companies_group.Post("/", middlewares.Protected(), companies.Create)
	companies_group.Get("/", middlewares.Protected(), middlewares.Pagination, companies.GetAll)
	companies_group.Get("/:id", middlewares.Protected(), companies.GetByID)
	companies_group.Put("/:id", middlewares.Protected(), companies.Update)
	companies_group.Delete("/:id", middlewares.Protected(), companies.Delete)

	// Company Users (nested)
	cu_group := api.Group("/companies/:company_id/users", middlewares.Protected())
	cu_group.Post("/", company_users.Add)
	cu_group.Get("/", middlewares.Pagination, company_users.List)
	cu_group.Get("/:user_id", company_users.Get)
	cu_group.Put("/:user_id", company_users.Update)
	cu_group.Delete("/:user_id", company_users.Remove)

	// Company Branches (private)
	branches_group := api.Group("/company-branches")
	branches_group.Post("/", middlewares.Protected(), company_branches.Create)
	branches_group.Get("/", middlewares.Protected(), middlewares.Pagination, company_branches.GetAll)
	branches_group.Get("/:id", middlewares.Protected(), company_branches.GetByID)
	branches_group.Put("/:id", middlewares.Protected(), company_branches.Update)
	branches_group.Delete("/:id", middlewares.Protected(), company_branches.Delete)

	// Divisions (private)
	divisions_group := api.Group("/divisions")
	divisions_group.Post("/", middlewares.Protected(), divisions.Create)
	divisions_group.Get("/", middlewares.Protected(), middlewares.Pagination, divisions.GetAll)
	divisions_group.Get("/:id", middlewares.Protected(), divisions.GetByID)
	divisions_group.Put("/:id", middlewares.Protected(), divisions.Update)
	divisions_group.Delete("/:id", middlewares.Protected(), divisions.Delete)

	// Audit Logs (private - admin only ideally, but we put it here)
	audit_group := api.Group("/audit-logs")
	audit_group.Get("/", middlewares.Protected(), middlewares.Pagination, audit_logs.GetAll)
	audit_group.Get("/:id", middlewares.Protected(), audit_logs.GetByID)

	// User Devices (private)
	device_group := api.Group("/user-devices")
	device_group.Get("/", middlewares.Protected(), user_devices.List)
	device_group.Post("/", middlewares.Protected(), user_devices.Register)
	device_group.Post("/:id/revoke", middlewares.Protected(), user_devices.Revoke)
	device_group.Patch("/:id/push-token", middlewares.Protected(), user_devices.UpdatePushToken)

	// 404
	app.Use(not_found)
}
