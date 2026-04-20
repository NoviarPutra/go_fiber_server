package divisions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func GetAll(c *fiber.Ctx) error {
	companyID := c.Query("company_id")
	if companyID == "" {
		return utils.BadRequest(c, "company_id wajib diisi")
	}

	page, ok := c.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}
	perPage, ok := c.Locals("per_page").(int)
	if !ok || perPage < 1 {
		perPage = 10
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewDivisionsService(db)

	ctx := c.Context()
	divisions, total, err := svc.GetAll(ctx, companyID, page, perPage)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessWithMeta(c, divisions, "Data divisi berhasil diambil", &types.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewDivisionsService(db)

	ctx := c.Context()
	division, err := svc.GetByID(ctx, id)
	if err != nil {
		if err == services.ErrDivisionNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, division, "Data divisi ditemukan")
}
