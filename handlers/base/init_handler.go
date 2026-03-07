package base

import "github.com/gofiber/fiber/v2"

func InitHandler(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "Welcome to the Go Fiber Server!",
		"code":    200,
		"data":    nil,
	})

}
