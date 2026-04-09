package integration

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/go_server/internal/middlewares"
)

func TestLoggerMiddleware_Detailed(t *testing.T) {
	// Set development mode to cover pretty JSON and colorized output
	os.Setenv("APP_ENV", "development")
	defer os.Setenv("APP_ENV", "testing")

	app := fiber.New()
	app.Use(middlewares.LoggerMiddleware())

	app.Get("/json", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status": "ok",
			"data": map[string]interface{}{
				"id":      123,
				"name":    "test",
				"active":  true,
				"nothing": nil,
				"list":    []int{1, 2, 3},
			},
		})
	})

	app.Get("/slow", func(c *fiber.Ctx) error {
		time.Sleep(150 * time.Millisecond) // Trigger snail/warning emoji
		return c.SendString("slow")
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(500, "intentional error")
	})

	app.Get("/301", func(c *fiber.Ctx) error {
		return c.Status(301).SendString("moved")
	})

	app.Get("/403", func(c *fiber.Ctx) error {
		return c.Status(403).SendString("forbidden")
	})

	app.Get("/429", func(c *fiber.Ctx) error {
		return c.Status(429).SendString("too many")
	})

	app.Get("/401", func(c *fiber.Ctx) error {
		return c.Status(401).SendString("unauthorized")
	})

	app.Get("/405", func(c *fiber.Ctx) error {
		return c.Status(405).SendString("not allowed")
	})

	t.Run("JSON_Logging", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/json", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Slow_Request_Logging", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/slow", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Error_Logging", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 500, resp.StatusCode)
	})

	t.Run("Other_Status_Logging", func(t *testing.T) {
		for _, path := range []string{"/301", "/403", "/429", "/401", "/405"} {
			req := httptest.NewRequest("GET", path, nil)
			_, _ = app.Test(req)
		}
	})
}
