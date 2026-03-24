package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/internal/types"
)

// SendResponse adalah helper internal — semua fungsi wajib lewat sini
func SendResponse[T any](c *fiber.Ctx, code int, success bool, msg string, data T, meta *types.Meta) error {
	return c.Status(code).JSON(types.StandardResponse[T]{
		Success: success,
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

// ─── Success Responses ────────────────────────────────────────────────────────

func Success[T any](c *fiber.Ctx, data T, msg string) error {
	return SendResponse(c, fiber.StatusOK, true, msg, data, nil)
}

func SuccessWithMeta[T any](c *fiber.Ctx, data T, msg string, meta *types.Meta) error {
	return SendResponse(c, fiber.StatusOK, true, msg, data, meta) // ✅ lewat SendResponse
}

func Created[T any](c *fiber.Ctx, data T, msg string) error {
	return SendResponse(c, fiber.StatusCreated, true, msg, data, nil)
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Error Responses ──────────────────────────────────────────────────────────

func ErrorResponse(c *fiber.Ctx, code int, msg string) error {
	return SendResponse[any](c, code, false, msg, nil, nil)
}

// Shorthand untuk error yang paling sering dipakai di handler
func BadRequest(c *fiber.Ctx, msg string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, msg)
}

func Unauthorized(c *fiber.Ctx, msg string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, msg)
}

func Forbidden(c *fiber.Ctx, msg string) error {
	return ErrorResponse(c, fiber.StatusForbidden, msg)
}

func NotFound(c *fiber.Ctx, msg string) error {
	return ErrorResponse(c, fiber.StatusNotFound, msg)
}

func InternalError(c *fiber.Ctx, msg string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, msg)
}
