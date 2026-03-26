package errors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGlobalErrorHandler(t *testing.T) {
	// Setup app dengan ErrorHandler kustom
	app := fiber.New(fiber.Config{
		ErrorHandler: GlobalErrorHandler,
	})

	t.Run("Handle_Fiber_Standard_Error_404", func(t *testing.T) {
		app.Get("/not-found", func(c *fiber.Ctx) error {
			return fiber.ErrNotFound // Memicu 404
		})

		req := httptest.NewRequest("GET", "/not-found", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 404, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		assert.Equal(t, false, result["success"])
		assert.Equal(t, "Not Found", result["message"])
	})

	t.Run("Handle_Generic_Internal_Error_500", func(t *testing.T) {
		app.Get("/panic", func(c *fiber.Ctx) error {
			return errors.New("database down")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		assert.Equal(t, "Terjadi kesalahan pada server", result["message"])
	})

	t.Run("Development_Mode_Stack_Trace_Logging", func(t *testing.T) {
		// Simulasikan environment development
		os.Setenv("APP_ENV", "development")
		defer os.Unsetenv("APP_ENV")

		app.Get("/stack", func(c *fiber.Ctx) error {
			// Gunakan pkg/errors untuk wrap stack trace
			return pkg_errors.WithStack(errors.New("critical failure"))
		})

		req := httptest.NewRequest("GET", "/stack", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		// Verifikasi manual log output jika perlu,
		// tapi di sini kita memastikan flow handler tidak pecah saat ada stack trace.
	})

	t.Run("Production_Mode_Silence_Stack_Trace", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		defer os.Unsetenv("APP_ENV")

		app.Get("/prod-error", func(c *fiber.Ctx) error {
			return pkg_errors.WithStack(errors.New("secret error"))
		})

		req := httptest.NewRequest("GET", "/prod-error", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		// Pastikan message tetap generik untuk user
		assert.Equal(t, "Terjadi kesalahan pada server", result["message"])
	})
}
