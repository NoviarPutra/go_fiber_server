package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal" // Sesuaikan dengan path bootstrap Anda
)

type RoutesTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *RoutesTestSuite) SetupSuite() {
	// Inisialisasi app secara penuh termasuk semua route dan middleware
	s.app = internal.Bootstrap(testDBPool)
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *RoutesTestSuite) TestPublicRoutes() {
	tests := []struct {
		name string
		path string
	}{
		{"Root_Path", "/"},
		{"Health_Check", "/health"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, _ := s.app.Test(req)

			// Minimal memastikan route terdaftar dan tidak 404
			s.Equal(200, resp.StatusCode)
		})
	}
}

func (s *RoutesTestSuite) TestAuthRoutes_Accessibility() {
	// Kita tes keberadaan endpoint, bukan logic (karena logic ada di auth_service_test)
	tests := []struct {
		name string
		path string
	}{
		{"Register_Endpoint", "/api/v1/auth/register"},
		{"Login_Endpoint", "/api/v1/auth/login"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// POST tanpa body akan memicu error 400/validation,
			// tapi membuktikan route tersebut TERDAFTAR (bukan 404).
			req := httptest.NewRequest("POST", tt.path, nil)
			resp, _ := s.app.Test(req)

			s.NotEqual(404, resp.StatusCode, "Route %s harusnya terdaftar", tt.path)
		})
	}
}

func (s *RoutesTestSuite) TestPrivateRoutes_Security() {
	s.Run("Users_GetAll_Should_Be_Protected", func() {
		// Mengetes apakah middleware.Protected() berfungsi di route ini
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		resp, _ := s.app.Test(req)

		// Harus 401 Unauthorized karena kita tidak mengirimkan Bearer Token
		s.Equal(401, resp.StatusCode, "Akses tanpa token ke route private harus gagal")
	})
}

func (s *RoutesTestSuite) TestNotFound_Handler() {
	s.Run("Custom_404_JSON_Response", func() {
		req := httptest.NewRequest("GET", "/api/v1/halaman-ngawur", nil)
		resp, _ := s.app.Test(req)

		s.Equal(404, resp.StatusCode)
		// Memastikan utils.NotFound bekerja (Response format JSON, bukan teks Fiber standar)
		s.Contains(resp.Header.Get("Content-Type"), "application/json")
	})
}

func TestRoutesIntegrity(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}
