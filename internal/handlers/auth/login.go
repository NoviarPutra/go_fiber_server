package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Login(c *fiber.Ctx) error {
	// 1. Parse body
	var req types.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Format request tidak valid")
	}

	// 2. Validasi input
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, format_validation_error(err))
	}

	// 3. Ambil DB dari context
	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewAuthService(db)

	// 4. Proses login
	result, err := svc.Login(c.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return utils.Unauthorized(c, "Email atau password salah")
		case errors.Is(err, services.ErrAccountInactive):
			return utils.Forbidden(c, "Akun Anda tidak aktif, hubungi administrator")
		default:
			return utils.InternalError(c, "Gagal login, coba lagi")
		}
	}

	// 5. Set cookies untuk FE (Security: HTTPOnly)
	utils.SetAuthCookies(c, result.AccessToken, result.RefreshToken)

	// Hilangkan token dari body agar tidak bisa dibaca JavaScript (Bullet-Proof)
	result.AccessToken = ""
	result.RefreshToken = ""

	return utils.Success(c, result, "Login berhasil")
}
