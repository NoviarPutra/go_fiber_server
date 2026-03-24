package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/go_server/internal"
)

func TestRouteNotFound_Integration(t *testing.T) {
	app := internal.Bootstrap(testDBPool)

	// Test route yang tidak ada untuk memastikan 404 handler Anda bekerja
	req := httptest.NewRequest("GET", "/api/v1/route-ngawur", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}
