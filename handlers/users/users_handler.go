package users

import "github.com/gofiber/fiber/v2"

func UsersHandler(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "Users endpoint",
		"code":    200,
		"data":    nil,
	})
}
