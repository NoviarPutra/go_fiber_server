package utils

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }() // Handle body close error

		is.Equal(fiber.StatusOK, resp.StatusCode)

		var body types.StandardResponse[map[string]interface{}]
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err, "Gagal decode JSON response")

		is.True(body.Success)
		is.Equal("Berhasil", body.Message)
		is.Equal(float64(1), body.Data["id"])
	})

	t.Run("Success With Meta", func(t *testing.T) {
		app.Get("/test-meta", func(c *fiber.Ctx) error {
			meta := &types.Meta{Total: 100, Page: 1}
			return SuccessWithMeta(c, []string{"data"}, "Berhasil", meta)
		})

		req := httptest.NewRequest("GET", "/test-meta", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		var body types.StandardResponse[[]string]
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		is.Equal(int64(100), int64(body.Meta.Total))
		is.Equal(int64(1), int64(body.Meta.Page))
	})

	t.Run("NoContent Response", func(t *testing.T) {
		app.Delete("/test-nocontent", func(c *fiber.Ctx) error {
			return NoContent(c)
		})

		req := httptest.NewRequest("DELETE", "/test-nocontent", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		is.Equal(fiber.StatusNoContent, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				app.Get(tt.route, tt.handler)
				req := httptest.NewRequest("GET", tt.route, nil)
				resp, err := app.Test(req)
				require.NoError(t, err)
				defer func() { _ = resp.Body.Close() }()

				is.Equal(tt.expectCode, resp.StatusCode)

				var body types.StandardResponse[any]
				err = json.NewDecoder(resp.Body).Decode(&body)
				require.NoError(t, err)
				is.False(body.Success)
				is.Equal(tt.expectMsg, body.Message)
			})
		}
	})
}
