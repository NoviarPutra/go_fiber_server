package user_devices

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Register(c *fiber.Ctx) error {
	var req types.RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Format data tidak valid")
	}

	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, format_validation_error(err))
	}

	userID := utils.GetUserIDFromCtx(c)
	if userID == "" {
		return utils.Unauthorized(c, "User tidak teridentifikasi")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewUserDevicesService(db)

	device, err := svc.RegisterDevice(c.Context(), userID, &req)
	if err != nil {
		return utils.InternalError(c, "Gagal mendaftarkan perangkat")
	}

	return utils.Success(c, device, "Perangkat berhasil didaftarkan")
}
