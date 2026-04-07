package companies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func GetAll(c *fiber.Ctx) error {
	page, ok := c.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}
	per_page, ok := c.Locals("per_page").(int)
	if !ok || per_page < 1 {
		per_page = 10
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompaniesService(db)

	companies, total, err := svc.GetAll(c.Context(), page, per_page)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data perusahaan")
	}

	return utils.SuccessWithMeta(c, companies, "Data perusahaan berhasil diambil", &types.Meta{
		Page:    page,
		PerPage: per_page,
		Total:   total,
	})
}

func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "ID diperlukan")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompaniesService(db)

	company, err := svc.GetByID(c.Context(), id)
	if err != nil {
		if err == services.ErrCompanyNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, company, "Data perusahaan ditemukan")
}
