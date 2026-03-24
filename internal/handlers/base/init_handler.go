package base

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/internal/utils"
)

func InitHandler(ctx *fiber.Ctx) error {
	return utils.Success[any](ctx, nil, "Welcome to the Go Fiber Server!")
}
