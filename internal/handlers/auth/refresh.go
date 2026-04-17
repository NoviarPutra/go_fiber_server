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
	_ = c.BodyParser(&req) // Ignore error as we also check cookies

	// 1. Ambil token dari body atau cookie
	refreshToken := req.RefreshToken
	if refreshToken == "" {
		refreshToken = c.Cookies(utils.CookieRefreshToken)
	}

	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token diperlukan")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	service := services.NewAuthService(db)

	resp, err := service.Refresh(c.Context(), refreshToken)
	if err != nil {
		if err == services.ErrRefreshTokenInvalid {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// 2. Update cookies baru
	utils.SetAuthCookies(c, resp.AccessToken, resp.RefreshToken)

	// Hilangkan token dari body untuk Web (mencegah XSS - dibaca JS).
	// Berikan lewat body jika origin adalah Mobile App.
	if c.Get("X-Client-Type") != "mobile" {
		resp.AccessToken = ""
		resp.RefreshToken = ""
	}

	return utils.Success(c, resp, "Token berhasil diperbarui")
}
