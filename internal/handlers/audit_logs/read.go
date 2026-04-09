package audit_logs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func GetAll(c *fiber.Ctx) error {
	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewAuditLogsService(db)

	page, ok := c.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}

	perPage, ok := c.Locals("per_page").(int)
	if !ok || perPage < 1 {
		perPage = 10
	}

	var query types.AuditLogQuery
	if err := c.QueryParser(&query); err != nil {
		return utils.BadRequest(c, "Parameter query tidak valid")
	}

	logs, total, err := svc.GetAll(c.Context(), query, page, perPage)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	meta := &types.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}

	return utils.SuccessWithMeta(c, logs, "Berhasil mengambil data audit log", meta)
}

func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewAuditLogsService(db)

	logData, err := svc.GetByID(c.Context(), id)
	if err != nil {
		if err == services.ErrAuditLogNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, logData, "Berhasil mengambil detail audit log")
}
