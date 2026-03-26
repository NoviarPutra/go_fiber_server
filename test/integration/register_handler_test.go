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

type RegisterHandlerTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *RegisterHandlerTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.app.Use(middlewares.DBMiddleware(testDBPool))
	s.app.Post("/api/v1/auth/register", auth.Register)
}

func (s *RegisterHandlerTestSuite) TearDownTest() {
	_, err := testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
	require.NoError(s.T(), err)
}

func (s *RegisterHandlerTestSuite) TestRegister_Success() {
	s.Run("Should_Create_User_With_Argon2_Hash", func() {
		payload := types.RegisterRequest{
			Email:    "newuser@officecore.id",
			Username: "newuser",
			Password: "StrongPassword123!",
		}
		body, err := json.Marshal(payload)
		require.NoError(s.T(), err)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, 10000)
		require.NoError(s.T(), err)
		defer func() { _ = resp.Body.Close() }()

		s.Equal(201, resp.StatusCode)

		var dbHash string
		err = testDBPool.QueryRow(context.Background(),
			"SELECT password_hash FROM users WHERE email = $1", payload.Email).Scan(&dbHash)

		s.NoError(err)
		s.Contains(dbHash, "$argon2id$", "Password harus disimpan dalam format Argon2id")
		s.NotEqual(payload.Password, dbHash, "Password tidak boleh disimpan sebagai plain text")
	})
}

func (s *RegisterHandlerTestSuite) TestRegister_Duplicate_Conflict() {
	existingEmail := "existing@officecore.id"
	existingUser := "existinguser"

	hash, err := utils.HashPassword("any-pass")
	require.NoError(s.T(), err)

	_, err = testDBPool.Exec(context.Background(),
		"INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)",
		existingEmail, existingUser, hash)
	require.NoError(s.T(), err)

	tests := []struct {
		name     string
		payload  types.RegisterRequest
		expected string
	}{
		{
			name: "Duplicate_Email",
			payload: types.RegisterRequest{
				Email:    existingEmail,
				Username: "differentuser",
				Password: "Password123!",
			},
			expected: "Email sudah terdaftar",
		},
		{
			name: "Duplicate_Username",
			payload: types.RegisterRequest{
				Email:    "newuser@officecore.id",
				Username: existingUser,
				Password: "Password123!",
			},
			expected: "Username sudah digunakan",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, err := json.Marshal(tt.payload)
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.app.Test(req, 10000)
			require.NoError(s.T(), err)
			defer func() { _ = resp.Body.Close() }()

			var res map[string]interface{}

			// FIX UTAMA: Tangkap error dari Decode
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(s.T(), err, "Gagal mendecode response error register")

			s.Equal(400, resp.StatusCode)
			s.Equal(tt.expected, res["message"])
		})
	}
}

func (s *RegisterHandlerTestSuite) TestRegister_Validation_Rules() {
	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"Invalid_Email", map[string]interface{}{"email": "bukan-email", "username": "user", "password": "pass"}},
		{"Password_Too_Short", map[string]interface{}{"email": "a@b.com", "username": "user", "password": "123"}},
		{"Empty_Fields", map[string]interface{}{"email": "", "username": "", "password": ""}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, err := json.Marshal(tt.payload)
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.app.Test(req)
			require.NoError(s.T(), err)
			defer func() { _ = resp.Body.Close() }()

			s.Equal(400, resp.StatusCode)
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	suite.Run(t, new(RegisterHandlerTestSuite))
}
