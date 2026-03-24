package middlewares

import (
	"errors"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/go_server/internal/utils"
)

func Protected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET tidak ditemukan di environment")
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			// FIX: Untuk HS256, cukup masukkan key []byte.
			// Library secara otomatis akan memvalidasi signature menggunakan HMAC.
			Key: []byte(secret),
		},
		ContextKey: "user_auth",
		// Jika ingin membatasi algoritma secara eksplisit di versi terbaru:
		// SigningKeys: map[string]jwtware.SigningKey{"HS256": {Key: []byte(secret)}},
		ErrorHandler: jwt_error_handler,
	})
}

func jwt_error_handler(c *fiber.Ctx, err error) error {
	// Filter error menggunakan errors.As atau errors.Is untuk akurasi

	// 1. Token hilang atau format "Bearer <token>" salah
	if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
		return utils.Unauthorized(c, "Token diperlukan atau format salah")
	}

	// 2. Token Expired
	if errors.Is(err, jwt.ErrTokenExpired) {
		return utils.Unauthorized(c, "Token kedaluwarsa, silakan login kembali")
	}

	// 3. Signature Mismatch (Inilah yang menangani manipulasi huruf belakang)
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return utils.Unauthorized(c, "Token tidak valid (Signature Mismatch)")
	}

	// 4. Fallback untuk error lainnya (Invalid claims, dll)
	return utils.Unauthorized(c, "Akses ditolak")
}
