package company_branches

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Create(c *fiber.Ctx) error {
	var req types.CreateCompanyBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Data tidak valid")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyBranchesService(db)

	branch, err := svc.Create(c.Context(), req)
	if err != nil {
		if err == services.ErrCompanyBranchNameExists {
			return utils.Conflict(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Created(c, branch, "Cabang perusahaan berhasil dibuat")
}
