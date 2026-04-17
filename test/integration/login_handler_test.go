package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/auth"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

type LoginHandlerTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *LoginHandlerTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.app.Use(middlewares.DBMiddleware(testDBPool))
	s.app.Post("/api/v1/auth/login", auth.Login)
}

func (s *LoginHandlerTestSuite) TearDownTest() {
	_, _ = testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
}

// ─── HELPERS ──────────────────────────────────────────────────────────────────

func (s *LoginHandlerTestSuite) seedUser(email, password string, active bool) {
	hash, err := utils.HashPassword(password)
	s.Require().NoError(err, "Gagal hashing password Argon2")

	username := email
	query := `INSERT INTO users (email, username, password_hash, is_active) 
              VALUES ($1, $2, $3, $4)`
	_, err = testDBPool.Exec(context.Background(), query, email, username, hash, active)
	s.Require().NoError(err, "Gagal seeding user ke database test")
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *LoginHandlerTestSuite) TestLogin_Success() {
	s.seedUser("ceo@officecore.id", "SecurePass123!", true)

	payload := types.LoginRequest{
		Email:    "ceo@officecore.id",
		Password: "SecurePass123!",
	}
	body, err := json.Marshal(payload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Argon2id butuh waktu, 10s timeout sudah tepat
	resp, err := s.app.Test(req, 10000)
	require.NoError(s.T(), err)
	defer func() { _ = resp.Body.Close() }() // Mencegah FD leak

	s.Equal(200, resp.StatusCode)

	var res map[string]interface{}
	// FIX: Tangkap error dari Decode untuk lolos errcheck
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(s.T(), err, "Gagal decode response login")

	s.Equal("Login berhasil", res["message"])

	data, ok := res["data"].(map[string]interface{})
	require.True(s.T(), ok, "Format data response tidak sesuai")
	s.Empty(data["access_token"], "JWT Access Token tidak boleh ada di body")
	s.Empty(data["refresh_token"], "JWT Refresh Token tidak boleh ada di body")

	// Verifikasi Cookie
	var accessTokenFound, refreshTokenFound bool
	for _, cookie := range resp.Cookies() {
		if cookie.Name == utils.CookieAccessToken {
			accessTokenFound = true
			s.NotEmpty(cookie.Value)
			s.True(cookie.HttpOnly)
		}
		if cookie.Name == utils.CookieRefreshToken {
			refreshTokenFound = true
			s.NotEmpty(cookie.Value)
			s.True(cookie.HttpOnly)
		}
	}
	s.True(accessTokenFound, "Access token harus di-set di HTTP-only cookie")
	s.True(refreshTokenFound, "Refresh token harus di-set di HTTP-only cookie")
}

func (s *LoginHandlerTestSuite) TestLogin_Security_Rejections() {
	s.seedUser("active@officecore.id", "password123", true)
	s.seedUser("inactive@officecore.id", "password123", false)

	tests := []struct {
		name       string
		email      string
		password   string
		expectCode int
	}{
		{"Wrong_Password", "active@officecore.id", "salah_pass", 401},
		{"Non_Existent_Email", "ghost@officecore.id", "password123", 401},
		{"Inactive_Account", "inactive@officecore.id", "password123", 403},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := types.LoginRequest{Email: tt.email, Password: tt.password}
			body, err := json.Marshal(payload)
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.app.Test(req, 10000)
			require.NoError(s.T(), err)
			defer func() { _ = resp.Body.Close() }()

			s.Equal(tt.expectCode, resp.StatusCode)
		})
	}
}

func (s *LoginHandlerTestSuite) TestLogin_Payload_Robustness() {
	s.Run("Malformed_JSON", func() {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString("{invalid:json}"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer func() { _ = resp.Body.Close() }()

		s.Equal(400, resp.StatusCode)
	})

	s.Run("Validation_Field_Missing", func() {
		payload := map[string]string{"email": "not-an-email"}
		body, err := json.Marshal(payload)
		require.NoError(s.T(), err)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer func() { _ = resp.Body.Close() }()

		s.Equal(400, resp.StatusCode)
	})
}

func TestLoginHandler(t *testing.T) {
	suite.Run(t, new(LoginHandlerTestSuite))
}
