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
	"github.com/yourusername/go_server/internal/utils" // Pastikan ada util Argon2 di sini
)

type LoginHandlerTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *LoginHandlerTestSuite) SetupSuite() {
	s.app = fiber.New()
	// Middleware DB wajib di-inject agar handler bisa ambil pool
	s.app.Use(middlewares.DBMiddleware(testDBPool))
	s.app.Post("/api/v1/auth/login", auth.Login)
}

func (s *LoginHandlerTestSuite) TearDownTest() {
	// Membersihkan data user setiap selesai satu test case agar tidak bentrok (Hermetic)
	_, _ = testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
}

// ─── HELPERS ──────────────────────────────────────────────────────────────────

func (s *LoginHandlerTestSuite) seedUser(email, password string, active bool) {
	hash, err := utils.HashPassword(password)
	s.Require().NoError(err, "Gagal hashing password Argon2")

	// FIX: Gunakan email prefix sebagai username agar UNIQUE constraint tidak jebol
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
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, 10000) // Argon2 butuh waktu lebih lama, beri timeout 10s
	s.NoError(err)
	s.Equal(200, resp.StatusCode)

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	s.Equal("Login berhasil", res["message"])
	data := res["data"].(map[string]interface{})
	s.NotEmpty(data["access_token"], "JWT Access Token tidak boleh kosong")
}

func (s *LoginHandlerTestSuite) TestLogin_Security_Rejections() {
	// Pindahkan seeding ke DALAM sub-test atau bersihkan secara eksplisit
	// Tapi cara terbaik adalah seeding satu kali di awal untuk semua skenario rejections ini:

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
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Tambahkan timeout 10s karena Argon2id sangat berat di CPU
			resp, err := s.app.Test(req, 10000)
			s.NoError(err)
			s.Equal(tt.expectCode, resp.StatusCode)
		})
	}
}

func (s *LoginHandlerTestSuite) TestLogin_Payload_Robustness() {
	s.Run("Malformed_JSON", func() {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString("{invalid:json}"))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := s.app.Test(req)
		s.Equal(400, resp.StatusCode)
	})

	s.Run("Validation_Field_Missing", func() {
		payload := map[string]string{"email": "not-an-email"} // Format email salah, password kosong
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := s.app.Test(req)
		s.Equal(400, resp.StatusCode)
	})
}

func TestLoginHandler(t *testing.T) {
	suite.Run(t, new(LoginHandlerTestSuite))
}
