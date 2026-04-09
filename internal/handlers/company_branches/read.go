package company_branches

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

	db, ok := c.Locals("db").(*pgxpool.Pool)
	if !ok {
		return utils.InternalError(c, "database connection not found")
	}
	svc := services.NewCompanyBranchesService(db)

	branches, totalCount, err := svc.GetAll(c.Context(), page, per_page)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessWithMeta(c, branches, "Berhasil mengambil data cabang perusahaan", &types.Meta{
		Page:    page,
		PerPage: per_page,
		Total:   totalCount,
	})
}

func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyBranchesService(db)

	branch, err := svc.GetByID(c.Context(), id)
	if err != nil {
		if err == services.ErrCompanyBranchNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, branch, "Berhasil mengambil data cabang perusahaan")
}
