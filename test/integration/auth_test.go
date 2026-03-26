package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal"
	"github.com/yourusername/go_server/internal/types"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	app     *fiber.App
	cleanup func()
	ctx     context.Context
}

func (s *AuthIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	dbPool, cleanup := SetupTestContainer(s.ctx)
	s.cleanup = cleanup
	s.app = internal.Bootstrap(dbPool)
}

// SetupTest berjalan SETIAP KALI sebelum satu s.Run(...)
// Ini memastikan email "test.user@officecore.id" tidak bentrok antar sub-test
func (s *AuthIntegrationTestSuite) SetupTest() {
	// Opsional: Tambahkan logika TRUNCATE users table di sini jika perlu
}

func (s *AuthIntegrationTestSuite) TearDownSuite() {
	s.cleanup()
}

func (s *AuthIntegrationTestSuite) TestAuthFlow() {
	// Gunakan data yang valid sesuai requirement logic Anda
	username := "testuser123"
	email := "test.user@officecore.id"
	password := "Secret123!"

	s.Run("Register_New_User", func() {
		payload := map[string]interface{}{
			"username": username, // TAMBAHKAN INI
			"email":    email,
			"password": password,
			"name":     "Test User",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Gunakan timeout yang aman untuk environment Docker/Colima
		resp, _ := s.app.Test(req, 30000)

		if resp.StatusCode != fiber.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			s.T().Errorf("Register Failed! Status: %d, Body: %s", resp.StatusCode, string(respBody))
		}

		s.Equal(fiber.StatusCreated, resp.StatusCode)
	})

	s.Run("Login_Success", func() {
		// Login biasanya menggunakan email atau username
		// Sesuaikan dengan logic handler login Anda
		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := s.app.Test(req, 30000)

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[map[string]interface{}]
		json.NewDecoder(resp.Body).Decode(&result)

		s.True(result.Success)
		s.NotEmpty(result.Data["access_token"])
	})
}

func TestAuthIntegration(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
