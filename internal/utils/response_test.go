package utils

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/go_server/internal/types"
)

func TestResponseSuite(t *testing.T) {
	app := fiber.New()
	is := assert.New(t)

	t.Run("Success Response", func(t *testing.T) {
		app.Get("/test-success", func(c *fiber.Ctx) error {
			return Success(c, fiber.Map{"id": 1}, "Berhasil")
		})

		req := httptest.NewRequest("GET", "/test-success", nil)
		resp, _ := app.Test(req)

		is.Equal(fiber.StatusOK, resp.StatusCode)

		var body types.StandardResponse[map[string]interface{}]
		json.NewDecoder(resp.Body).Decode(&body)

		is.True(body.Success)
		is.Equal("Berhasil", body.Message)
		is.Equal(float64(1), body.Data["id"]) // JSON numbers are float64 in Go
	})

	t.Run("Success With Meta", func(t *testing.T) {
		app.Get("/test-meta", func(c *fiber.Ctx) error {
			meta := &types.Meta{Total: 100, Page: 1}
			return SuccessWithMeta(c, []string{"data"}, "Berhasil", meta)
		})

		req := httptest.NewRequest("GET", "/test-meta", nil)
		resp, _ := app.Test(req)

		var body types.StandardResponse[[]string]
		json.NewDecoder(resp.Body).Decode(&body)

		// FIX: Cast ke int64 agar match dengan hasil unmarshal JSON angka
		is.Equal(int64(100), int64(body.Meta.Total))
		is.Equal(int64(1), int64(body.Meta.Page))
	})

	t.Run("Created Response", func(t *testing.T) {
		app.Post("/test-created", func(c *fiber.Ctx) error {
			return Created(c, fiber.Map{"id": 99}, "User dibuat")
		})

		req := httptest.NewRequest("POST", "/test-created", nil)
		resp, _ := app.Test(req)

		is.Equal(fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("NoContent Response", func(t *testing.T) {
		app.Delete("/test-nocontent", func(c *fiber.Ctx) error {
			return NoContent(c)
		})

		req := httptest.NewRequest("DELETE", "/test-nocontent", nil)
		resp, _ := app.Test(req)

		is.Equal(fiber.StatusNoContent, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		is.Empty(body)
	})

	t.Run("Error Shorthands", func(t *testing.T) {
		tests := []struct {
			name       string
			route      string
			handler    func(*fiber.Ctx) error
			expectCode int
			expectMsg  string
		}{
			{"BadRequest", "/400", func(c *fiber.Ctx) error { return BadRequest(c, "Bad") }, 400, "Bad"},
			{"Unauthorized", "/401", func(c *fiber.Ctx) error { return Unauthorized(c, "Unauth") }, 401, "Unauth"},
			{"Forbidden", "/403", func(c *fiber.Ctx) error { return Forbidden(c, "Forb") }, 403, "Forb"},
			{"NotFound", "/404", func(c *fiber.Ctx) error { return NotFound(c, "NotF") }, 404, "NotF"},
			{"InternalError", "/500", func(c *fiber.Ctx) error { return InternalError(c, "Err") }, 500, "Err"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				app.Get(tt.route, tt.handler)
				req := httptest.NewRequest("GET", tt.route, nil)
				resp, _ := app.Test(req)

				is.Equal(tt.expectCode, resp.StatusCode)

				var body types.StandardResponse[any]
				json.NewDecoder(resp.Body).Decode(&body)
				is.False(body.Success)
				is.Equal(tt.expectMsg, body.Message)
			})
		}
	})
}
