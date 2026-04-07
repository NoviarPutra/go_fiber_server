package companies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/utils"
)

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequest(c, "ID diperlukan")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompaniesService(db)

	if err := svc.Delete(c.Context(), id); err != nil {
		if err == services.ErrCompanyNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success[any](c, nil, "Data perusahaan berhasil dihapus")
}
