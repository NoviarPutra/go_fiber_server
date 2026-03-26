package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTTestSuite(t *testing.T) {
	// Buat objek assert agar pemanggilan fungsi lebih ringkas
	is := assert.New(t)

	const testSecret = "my-super-secret-key-123"
	os.Setenv("JWT_SECRET", testSecret)

	userID := "user-888"
	email := "budiawan@fedora.local"

	t.Run("Generate & Parse Access Token Success", func(t *testing.T) {
		tokenStr, err := GenerateAccessToken(userID, email)
		is.NoError(err)
		is.NotEmpty(tokenStr)

		claims, err := ParseClaims(tokenStr)
		is.NoError(err)
		is.Equal(userID, claims.UserID)
		is.Equal(email, claims.Email)
		is.Equal(userID, claims.Subject)
	})

	t.Run("Generate Refresh Token Success", func(t *testing.T) {
		tokenStr, expiry, err := GenerateRefreshToken(userID)
		is.NoError(err)
		is.NotEmpty(tokenStr)
		is.WithinDuration(time.Now().Add(7*24*time.Hour), expiry, 5*time.Second)
	})

	t.Run("Fail when JWT_SECRET is Missing", func(t *testing.T) {
		// Gunakan t.Setenv untuk isolasi yang lebih baik di Go modern
		t.Setenv("JWT_SECRET", "")

		token, err := GenerateAccessToken(userID, email)
		is.Error(err)
		is.Empty(token)
		is.Contains(err.Error(), "tidak dikonfigurasi")
	})

	t.Run("Parse Invalid Token Format", func(t *testing.T) {
		claims, err := ParseClaims("ini.bukan.token")
		is.Error(err)
		is.Nil(claims) // Sekarang argumennya sudah benar karena menggunakan objek 'is'
	})

	t.Run("Parse Expired Token", func(t *testing.T) {
		secret := []byte(testSecret)
		claims := &JWTClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, _ := token.SignedString(secret)

		parsedClaims, err := ParseClaims(expiredToken)
		is.Error(err)
		is.Nil(parsedClaims)
		is.Contains(err.Error(), "token is expired")
	})

	t.Run("Parse Token with Wrong Signing Method", func(t *testing.T) {
		claims := &JWTClaims{UserID: userID}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		unsignedToken, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

		parsedClaims, err := ParseClaims(unsignedToken)
		is.Error(err)
		is.Nil(parsedClaims)
	})
}
