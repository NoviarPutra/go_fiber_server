package integration

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/go_server/internal"
)

func TestHealthHandler_Integration(t *testing.T) {
	// 1. Inisialisasi App Fiber menggunakan pool database dari Testcontainers
	// Karena Bootstrap(db) menyuntikkan middleware DBMiddleware(db),
	// c.Locals("db") akan terisi secara otomatis.
	app := internal.Bootstrap(testDBPool)

	// 2. Buat Request ke endpoint /health
	req := httptest.NewRequest("GET", "/health", nil)

	// 3. Eksekusi Request
	resp, err := app.Test(req, -1) // -1 menghilangkan limit timeout default fiber test

	// 4. Assertions (Standar Enterprise)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 5. Validasi Struktur Response Body
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	// Cek apakah success true (asumsi utils.Success mengembalikan field 'success')
	assert.True(t, body["success"].(bool))

	// Cek data stats database
	data := body["data"].(map[string]interface{})
	assert.NotNil(t, data["db_pool_total"])
	assert.NotNil(t, data["db_pool_idle"])

	assert.Equal(t, "Server dan database berjalan normal", body["message"])
}
