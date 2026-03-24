package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Constants untuk keamanan
const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

func Pagination(c *fiber.Ctx) error {
	// 1. Ambil query dengan penanganan error yang benar
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil || perPage < 1 {
		perPage = DefaultPerPage
	}

	// 2. Batasi MaxPerPage agar database tidak kewalahan (proteksi DDoS)
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	// 3. Simpan ke Locals
	c.Locals("page", page)
	c.Locals("per_page", perPage)

	return c.Next()
}
