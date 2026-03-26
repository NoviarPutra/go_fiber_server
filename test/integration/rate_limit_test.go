package integration

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type RateLimitTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *RateLimitTestSuite) SetupTest() {
	s.app = fiber.New()
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *RateLimitTestSuite) TestRateLimit_UnderLimit() {
	// 1. Setup Dev Mode (Limit 1000)
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	s.app.Get("/test", middlewares.RateLimitMiddleware(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// 2. Kirim beberapa request
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := s.app.Test(req)
		s.Equal(200, resp.StatusCode, "Request ke-%d harusnya sukses", i+1)
	}
}

func (s *RateLimitTestSuite) TestRateLimit_ExceedLimit() {
	// 1. Paksa mode produksi agar limit kecil (100) atau buat test limit sendiri
	// Agar test cepat, kita bisa memanipulasi env jika middleware membacanya saat inisialisasi
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	s.app.Get("/limited", middlewares.RateLimitMiddleware(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// 2. Habiskan kuota (Max: 100)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/limited", nil)
		resp, _ := s.app.Test(req)
		s.Require().Equal(200, resp.StatusCode)
	}

	// 3. Request ke-101 HARUS gagal (429)
	req := httptest.NewRequest("GET", "/limited", nil)
	resp, _ := s.app.Test(req)

	s.Equal(429, resp.StatusCode, "Harusnya mengembalikan 429 Too Many Requests")
}

func (s *RateLimitTestSuite) TestRateLimit_CustomError_Format() {
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	s.app.Get("/error-format", middlewares.RateLimitMiddleware(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Langsung habiskan limit
	for i := 0; i < 101; i++ {
		req := httptest.NewRequest("GET", "/error-format", nil)
		resp, _ := s.app.Test(req)

		if i == 100 {
			// Cek apakah response body menggunakan utils.ErrorResponse (JSON)
			// bukan plain text default dari Fiber
			s.Equal(429, resp.StatusCode)
			s.Contains(resp.Header.Get("Content-Type"), "application/json")
		}
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
