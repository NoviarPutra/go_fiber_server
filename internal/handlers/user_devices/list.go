package user_devices

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/utils"
)

func List(c *fiber.Ctx) error {
	userID := utils.GetUserIDFromCtx(c)
	if userID == "" {
		return utils.Unauthorized(c, "User tidak teridentifikasi")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewUserDevicesService(db)

	devices, err := svc.ListDevices(c.Context(), userID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil daftar perangkat")
	}

	return utils.Success(c, devices, "Daftar perangkat berhasil diambil")
}
