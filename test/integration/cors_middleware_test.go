package integration

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type CorsMiddlewareTestSuite struct {
	suite.Suite
}

func (s *CorsMiddlewareTestSuite) TestCors_Development_Mode() {
	// 1. Setup Environment
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	app := fiber.New()
	app.Use(middlewares.CorsMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	// 2. Action: Simulasikan Preflight Request (OPTIONS) atau GET biasa
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000") // Simulate frontend origin

	resp, _ := app.Test(req)

	// 3. Assert: Di dev mode, harus mengembalikan wildcard (*)
	s.Equal(200, resp.StatusCode)
	s.Equal("*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func (s *CorsMiddlewareTestSuite) TestCors_Production_Mode() {
	// 1. Setup Environment Produksi dengan origin spesifik
	os.Setenv("APP_ENV", "production")
	os.Setenv("ALLOWED_ORIGINS", "https://hadir.officecore.id")
	defer os.Unsetenv("APP_ENV")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	app := fiber.New()
	app.Use(middlewares.CorsMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://hadir.officecore.id")

	resp, _ := app.Test(req)

	// 2. Assert: Harus mengembalikan origin yang diizinkan, bukan wildcard
	s.Equal("https://hadir.officecore.id", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCorsMiddleware(t *testing.T) {
	suite.Run(t, new(CorsMiddlewareTestSuite))
}
