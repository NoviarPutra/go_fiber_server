package health

import "github.com/gofiber/fiber/v2"

func HealthHandler(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is healthy",
		"code":    200,
		"data":    nil,
	})
}
