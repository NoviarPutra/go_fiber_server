package utils

import (
	"fmt"
	"os"
	"time"

	github_com_gofiber_fiber_v2 "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func get_jwt_secret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET tidak dikonfigurasi")
	}
	return []byte(secret), nil
}

func GenerateAccessToken(user_id, email string) (string, error) {
	secret, err := get_jwt_secret()
	if err != nil {
		return "", err
	}

	claims := JWTClaims{
		UserID: user_id,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user_id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(user_id string) (string, time.Time, error) {
	secret, err := get_jwt_secret()
	if err != nil {
		return "", time.Time{}, err
	}

	expires_at := time.Now().Add(RefreshTokenExpiry)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expires_at),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user_id,
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	return signed, expires_at, err
}

// ParseClaims — dipakai middleware atau handler untuk ambil data user dari token
func ParseClaims(token_string string) (*JWTClaims, error) {
	secret, err := get_jwt_secret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(token_string, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method tidak valid")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}

	return claims, nil
}

// GetUserIDFromCtx mengambil user ID dari token yang di set JWT Middleware
func GetUserIDFromCtx(c *github_com_gofiber_fiber_v2.Ctx) string {
	userToken, ok := c.Locals("user_auth").(*jwt.Token)
	if !ok {
		return ""
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	userID, _ := claims["user_id"].(string)
	return userID
}
