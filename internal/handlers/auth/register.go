package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Register(c *fiber.Ctx) error {
	// 1. Parse body
	var req types.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Format request tidak valid")
	}

	// 2. Validasi — pakai validate & format_validation_error dari validator.go
	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(c, format_validation_error(err))
	}

	// 3. Ambil DB dari context
	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewAuthService(db)

	// 4. Proses registrasi
	user, err := svc.Register(c.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists):
			return utils.BadRequest(c, "Email sudah terdaftar")
		case errors.Is(err, services.ErrUsernameAlreadyExists):
			return utils.BadRequest(c, "Username sudah digunakan")
		default:
			return utils.InternalError(c, "Gagal membuat akun, coba lagi")
		}
	}

	return utils.Created(c, user, "Akun berhasil dibuat")
}
