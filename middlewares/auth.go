package middlewares

import (
	"os"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "sangat-rahasia"
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		// Future-proof: Definisikan secara eksplisit di mana data user disimpan
		ContextKey: "user_auth",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Cek apakah token memang tidak ada atau cuma salah format
			message := "Akses ditolak: Token tidak valid"
			if strings.Contains(err.Error(), "exp") {
				message = "Sesi Anda telah berakhir, silakan login kembali"
			} else if strings.Contains(err.Error(), "missing") {
				message = "Token diperlukan untuk mengakses resource ini"
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": message,
				// Tambahkan detail error jika di mode development saja
				"debug": err.Error(),
			})
		},
	})
}
