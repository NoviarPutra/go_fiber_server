package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Revoke(c *fiber.Ctx) error {
	var req types.RevokeTokenRequest
	_ = c.BodyParser(&req)

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		refreshToken = c.Cookies(utils.CookieRefreshToken)
	}

	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh token diperlukan")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	service := services.NewAuthService(db)

	err := service.Revoke(c.Context(), refreshToken)
	if err != nil {
		if err == services.ErrRefreshTokenInvalid {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// 2. Bersihkan cookies
	utils.ClearAuthCookies(c)

	return utils.Success[any](c, nil, "Token berhasil dicabut (Logout)")
}
