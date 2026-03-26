package integration

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type RecoveryTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *RecoveryTestSuite) SetupTest() {
	s.app = fiber.New()
	s.app.Use(middlewares.RecoveryMiddleware())
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *RecoveryTestSuite) TestRecovery_Panic_Handling() {
	s.Run("Should_Recover_From_Panic_And_Return_500", func() {
		// 1. Buat route yang sengaja melakukan panic
		s.app.Get("/boom", func(c *fiber.Ctx) error {
			panic("something went terribly wrong")
		})

		// 2. Kirim request
		req := httptest.NewRequest("GET", "/boom", nil)
		resp, err := s.app.Test(req)

		// 3. Assert: App tidak boleh crash (err == nil) dan status harus 500
		s.NoError(err, "Server harusnya tetap hidup setelah panic")
		s.Equal(500, resp.StatusCode, "Status code harus 500 Internal Server Error")
	})
}

func (s *RecoveryTestSuite) TestRecovery_Environment_Logic() {
	s.Run("Development_Mode_Config", func() {
		os.Setenv("APP_ENV", "development")
		defer os.Unsetenv("APP_ENV")

		// Re-init middleware untuk membaca env baru
		app := fiber.New()
		app.Use(middlewares.RecoveryMiddleware())

		app.Get("/panic-dev", func(c *fiber.Ctx) error {
			panic("dev panic")
		})

		req := httptest.NewRequest("GET", "/panic-dev", nil)
		resp, _ := app.Test(req)

		s.Equal(500, resp.StatusCode)
		// Secara internal stack trace tereksekusi, server tetap aman.
	})

	s.Run("Production_Mode_Config", func() {
		os.Setenv("APP_ENV", "production")
		defer os.Unsetenv("APP_ENV")

		app := fiber.New()
		app.Use(middlewares.RecoveryMiddleware())

		app.Get("/panic-prod", func(c *fiber.Ctx) error {
			panic("prod panic")
		})

		req := httptest.NewRequest("GET", "/panic-prod", nil)
		resp, _ := app.Test(req)

		s.Equal(500, resp.StatusCode)
	})
}

func (s *RecoveryTestSuite) TestRecovery_Nil_Pointer_Safety() {
	s.Run("Should_Recover_From_Runtime_Panic", func() {
		s.app.Get("/trigger-panic", func(c *fiber.Ctx) error {
			if c.Method() != "" {
				panic("runtime error: manual trigger")
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/trigger-panic", nil)

		// Perbaikan di sini: tangkap response dan error secara terpisah
		resp, err := s.app.Test(req, -1)

		s.NoError(err, "App tidak boleh crash setelah panic")
		s.Equal(500, resp.StatusCode)

		// Opsional: Pastikan body tetap terbaca jika perlu
		defer resp.Body.Close()
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	suite.Run(t, new(RecoveryTestSuite))
}
