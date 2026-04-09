package company_branches

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/utils"
)

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyBranchesService(db)

	err := svc.Delete(c.Context(), id)
	if err != nil {
		if err == services.ErrCompanyBranchNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.NoContent(c)
}
