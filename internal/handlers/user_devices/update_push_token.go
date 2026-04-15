package user_devices

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func UpdatePushToken(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	if deviceID == "" {
		return utils.BadRequest(c, "ID perangkat diperlukan")
	}

	var req types.UpdatePushTokenRequest
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

	err := svc.UpdatePushToken(c.Context(), userID, deviceID, &req)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			return utils.NotFound(c, "Perangkat tidak ditemukan atau sudah tidak aktif")
		}
		return utils.InternalError(c, "Gagal memperbarui push token")
	}

	return utils.Success[any](c, nil, "Push token berhasil diperbarui")
}
