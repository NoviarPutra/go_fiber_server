package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/utils"
)

func HealthHandler(c *fiber.Ctx) error {
	db := c.Locals("db").(*pgxpool.Pool)

	ping_ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := db.Ping(ping_ctx); err != nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Database tidak merespons")
	}

	stats := db.Stat()
	return utils.Success(c, fiber.Map{
		"db_pool_total": stats.TotalConns(),
		"db_pool_idle":  stats.IdleConns(),
	}, "Server dan database berjalan normal")
}
