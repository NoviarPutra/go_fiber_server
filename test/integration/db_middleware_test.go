package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type DBMiddlewareTestSuite struct {
	suite.Suite
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *DBMiddlewareTestSuite) TestDBMiddleware_Success_Injection() {
	// 1. Setup App & Middleware
	// Kita gunakan testDBPool yang sudah di-init oleh TestMain
	s.Require().NotNil(testDBPool, "testDBPool harus sudah terinisialisasi")

	app := fiber.New()
	app.Use(middlewares.DBMiddleware(testDBPool))

	// 2. Route untuk verifikasi apakah DB ada di Locals
	app.Get("/verify-db", func(c *fiber.Ctx) error {
		db := c.Locals("db")

		// Assert tipe datanya benar
		pool, ok := db.(*pgxpool.Pool)
		if !ok || pool == nil {
			return c.Status(500).SendString("DB not injected correctly")
		}

		return c.SendString("DB Injected")
	})

	// 3. Action
	req := httptest.NewRequest("GET", "/verify-db", nil)
	resp, _ := app.Test(req)

	// 4. Assert
	s.Equal(200, resp.StatusCode)
}

func (s *DBMiddlewareTestSuite) TestDBMiddleware_Panic_On_Nil() {
	// Mengetes "Guard" logic: func DBMiddleware(db *pgxpool.Pool)
	// Kita harus memastikan aplikasi panic jika diberikan nil pool

	s.Run("Should_Panic_When_Pool_Is_Nil", func() {
		defer func() {
			r := recover()
			s.NotNil(r, "Middleware harusnya panic jika db nil")
			s.Contains(r, "database pool tidak boleh nil")
		}()

		// Ini akan memicu panic
		_ = middlewares.DBMiddleware(nil)
	})
}

func TestDBMiddleware(t *testing.T) {
	suite.Run(t, new(DBMiddlewareTestSuite))
}
