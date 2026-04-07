package companies

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
	if id == "" {
		return utils.BadRequest(c, "ID diperlukan")
	}

	var req types.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Data tidak valid")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompaniesService(db)

	company, err := svc.Update(c.Context(), id, req)
	if err != nil {
		if err == services.ErrCompanyNotFound {
			return utils.NotFound(c, err.Error())
		}
		if err == services.ErrCompanyCodeExists {
			return utils.BadRequest(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, company, "Data perusahaan berhasil diperbarui")
}
