package divisions

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Create(c *fiber.Ctx) error {
	var req types.CreateDivisionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Data tidak valid")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewDivisionsService(db)

	ctx := utils.InjectAuditInfo(c, utils.GetUserIDFromCtx(c), "")
	division, err := svc.Create(ctx, req)
	if err != nil {
		if err == services.ErrDivisionCodeExists {
			return utils.Conflict(c, err.Error())
		}
		if err == services.ErrCompanyNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Created(c, division, "Divisi berhasil dibuat")
}
