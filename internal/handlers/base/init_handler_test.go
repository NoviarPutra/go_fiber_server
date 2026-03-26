package base

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestInitHandler(t *testing.T) {
	// 1. Setup Fiber App
	app := fiber.New()
	app.Get("/", InitHandler)

	t.Run("Should_Return_200_With_Welcome_Message", func(t *testing.T) {
		// 2. Create Request
		req := httptest.NewRequest("GET", "/", nil)

		// 3. Execute Request
		resp, err := app.Test(req)

		// 4. Assertions
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		// 5. Verify JSON Body Structure
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)

		assert.NoError(t, err)
		assert.Equal(t, true, result["success"], "Flag success harus true")
		assert.Equal(t, "Welcome to the Go Fiber Server!", result["message"])
		assert.Nil(t, result["data"], "Data harus nil sesuai implementasi handler")
	})

	t.Run("Should_Handle_Method_Not_Allowed", func(t *testing.T) {
		// Mengetes ketangguhan routing (Edge Case)
		req := httptest.NewRequest("POST", "/", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 405, resp.StatusCode, "Endpoint ini hanya boleh menerima GET")
	})
}
