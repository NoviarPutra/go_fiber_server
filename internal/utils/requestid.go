package utils

import "github.com/gofiber/fiber/v2"

func GetRequestID(c *fiber.Ctx) string {
	val := c.Locals("requestid")
	if id, ok := val.(string); ok {
		return id
	}
	return "unknown"
}
