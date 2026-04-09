package company_branches

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req types.UpdateCompanyBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Data tidak valid")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyBranchesService(db)

	companyID := ""
	if req.CompanyID != nil {
		companyID = *req.CompanyID
	}
	ctx := utils.InjectAuditInfo(c, utils.GetUserIDFromCtx(c), companyID)
	branch, err := svc.Update(ctx, id, req)
	if err != nil {
		if err == services.ErrCompanyBranchNotFound {
			return utils.NotFound(c, err.Error())
		}
		if err == services.ErrCompanyBranchNameExists {
			return utils.Conflict(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, branch, "Cabang perusahaan berhasil diperbarui")
}
