package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBMiddleware(db *pgxpool.Pool) fiber.Handler {
	// Guard: pastikan db tidak nil sebelum inject
	if db == nil {
		panic("DBMiddleware: database pool tidak boleh nil")
	}

	return func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	}
}
