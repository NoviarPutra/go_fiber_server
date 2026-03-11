package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/config"
	"github.com/yourusername/go_server/utils"
)

func HealthHandler(ctx *fiber.Ctx) error {
	err := config.DB.Ping(ctx.Context())
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusServiceUnavailable, "Database connection lost")
	}
	return utils.Success[any](ctx, nil, "Server and Database are healthy")
}
