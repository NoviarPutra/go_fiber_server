package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func GetAll(c *fiber.Ctx) error {
	// 1. Ambil pagination dari context (di-set oleh Pagination middleware)
	page, ok := c.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}
	per_page, ok := c.Locals("per_page").(int)
	if !ok || per_page < 1 {
		per_page = 10
	}

	// 2. Ambil DB dari context
	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewUsersService(db)

	// 3. Ambil data dari service
	users, total, err := svc.GetUsers(c.Context(), page, per_page)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data users")
	}

	// 4. Response dengan meta pagination
	return utils.SuccessWithMeta(c, users, "Data users berhasil diambil", &types.Meta{
		Page:      page,
		PerPage:   per_page,
		Total:     total,
		RequestID: utils.GetRequestID(c),
	})
}
