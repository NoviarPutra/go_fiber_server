package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
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
	// Inject DB Pool agar c.Locals("db") tersedia
	s.app.Use(middlewares.DBMiddleware(testDBPool))
	s.app.Post("/api/v1/auth/register", auth.Register)
}

func (s *RegisterHandlerTestSuite) TearDownTest() {
	// Bersihkan tabel users setelah setiap test case
	_, _ = testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *RegisterHandlerTestSuite) TestRegister_Success() {
	s.Run("Should_Create_User_With_Argon2_Hash", func() {
		payload := types.RegisterRequest{
			Email:    "newuser@officecore.id",
			Username: "newuser",
			Password: "StrongPassword123!",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Argon2id butuh waktu, gunakan timeout yang cukup
		resp, err := s.app.Test(req, 10000)
		s.NoError(err)
		s.Equal(201, resp.StatusCode)

		// Verifikasi Database: Pastikan password di-hash (Argon2 biasanya dimulai dengan $argon2id$)
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

	hash, _ := utils.HashPassword("any-pass")
	_, _ = testDBPool.Exec(context.Background(),
		"INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)",
		existingEmail, existingUser, hash)

	tests := []struct {
		name     string
		payload  types.RegisterRequest
		expected string
	}{
		{
			name: "Duplicate_Email",
			payload: types.RegisterRequest{
				Email:    existingEmail,
				Username: "differentuser", // FIX: Hapus dash agar lolos validasi alphanum
				Password: "Password123!",
			},
			expected: "Email sudah terdaftar",
		},
		{
			name: "Duplicate_Username",
			payload: types.RegisterRequest{
				Email:    "newuser@officecore.id", // FIX: Email unik agar tidak kena error email duluan
				Username: existingUser,
				Password: "Password123!",
			},
			expected: "Username sudah digunakan",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := s.app.Test(req, 10000)

			var res map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&res)

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
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := s.app.Test(req)
			s.Equal(400, resp.StatusCode)
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	suite.Run(t, new(RegisterHandlerTestSuite))
}
