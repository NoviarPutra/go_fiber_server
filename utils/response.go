package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/types"
)

// SendResponse adalah helper internal
func SendResponse[T any](c *fiber.Ctx, code int, success bool, msg string, data T, meta *types.Meta) error {
	return c.Status(code).JSON(types.StandardResponse[T]{
		Success: success,
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

func SuccessWithMeta[T any](c *fiber.Ctx, data T, msg string, meta *types.Meta) error {
	return c.Status(fiber.StatusOK).JSON(types.StandardResponse[T]{
		Success: true,
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

// Success (200)
func Success[T any](c *fiber.Ctx, data T, msg string) error {
	return SendResponse(c, fiber.StatusOK, true, msg, data, nil)
}

// Created (201)
func Created[T any](c *fiber.Ctx, data T, msg string) error {
	return SendResponse(c, fiber.StatusCreated, true, msg, data, nil)
}

// NoContent (204)
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func ErrorResponse(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(types.StandardResponse[any]{
		Success: false,
		Message: msg,
		Data:    nil,
	})
}
