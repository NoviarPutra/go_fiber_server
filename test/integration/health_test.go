package integration

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/go_server/internal"
)

func TestHealthHandler_Integration(t *testing.T) {
	app := internal.Bootstrap(testDBPool)

	t.Run("Success Path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)

		// Pastikan data stats DB muncul (ini menaikkan coverage di handler)
		data := body["data"].(map[string]interface{})
		assert.NotNil(t, data["db_pool_total"])
	})

	t.Run("404 Route Path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/invalid-route", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 404, resp.StatusCode)
	})
}
