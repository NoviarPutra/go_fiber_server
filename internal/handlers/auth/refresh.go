package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Refresh(c *fiber.Ctx) error {
	var req types.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	if err := validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, format_validation_error(err))
	}

	db := c.Locals("db").(*pgxpool.Pool)
	service := services.NewAuthService(db)

	resp, err := service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		if err == services.ErrRefreshTokenInvalid {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, resp, "Token berhasil diperbarui")
}
