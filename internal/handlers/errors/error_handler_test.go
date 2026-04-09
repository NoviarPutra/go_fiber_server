package errors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: GlobalErrorHandler,
	})

	t.Run("Handle_Fiber_Standard_Error_404", func(t *testing.T) {
		app.Get("/not-found", func(c *fiber.Ctx) error {
			return fiber.ErrNotFound
		})

		req := httptest.NewRequest("GET", "/not-found", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, 404, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "Gagal unmarshal response body")

		assert.Equal(t, false, result["success"])
		assert.Equal(t, "Not Found", result["message"])
	})

	t.Run("Handle_Generic_Internal_Error_500", func(t *testing.T) {
		app.Get("/panic", func(c *fiber.Ctx) error {
			return errors.New("database down")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, 500, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Terjadi kesalahan pada server", result["message"])
	})
}
