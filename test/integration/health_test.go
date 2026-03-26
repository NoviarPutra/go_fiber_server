package integration

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/go_server/internal"
)

func TestHealthHandler_Integration(t *testing.T) {
	app := internal.Bootstrap(testDBPool)

	t.Run("Success Path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)

		// FIX: Tambahkan pengecekan error fundamental dan penutupan body
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, 200, resp.StatusCode)

		var body map[string]interface{}

		// FIX UTAMA: Tangkap dan periksa error return value dari Decode
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err, "Gagal mendecode JSON health response")

		// Pastikan data stats DB muncul
		data, ok := body["data"].(map[string]interface{})
		require.True(t, ok, "Format field 'data' tidak sesuai")
		assert.NotNil(t, data["db_pool_total"])
	})

	t.Run("404 Route Path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/invalid-route", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, 404, resp.StatusCode)
	})
}
