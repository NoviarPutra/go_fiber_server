package user_devices

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/utils"
)

func Revoke(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	if deviceID == "" {
		return utils.BadRequest(c, "ID perangkat diperlukan")
	}

	userID := utils.GetUserIDFromCtx(c)
	if userID == "" {
		return utils.Unauthorized(c, "User tidak teridentifikasi")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewUserDevicesService(db)

	err := svc.RevokeDevice(c.Context(), userID, deviceID)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			return utils.NotFound(c, "Perangkat tidak ditemukan atau sudah tidak aktif")
		}
		return utils.InternalError(c, "Gagal menonaktifkan perangkat")
	}

	return utils.Success[any](c, nil, "Perangkat berhasil dinonaktifkan")
}
