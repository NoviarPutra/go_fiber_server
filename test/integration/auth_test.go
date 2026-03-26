package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
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

func (s *AuthIntegrationTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *AuthIntegrationTestSuite) TestAuthFlow() {
	username := "testuser123"
	email := "test.user@officecore.id"
	password := "Secret123!"

	s.Run("Register_New_User", func() {
		payload := map[string]interface{}{
			"username": username,
			"email":    email,
			"password": password,
			"name":     "Test User",
		}
		body, err := json.Marshal(payload)
		require.NoError(s.T(), err)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, 30000)
		require.NoError(s.T(), err)
		defer func() { _ = resp.Body.Close() }() // Bullet-proof: Tutup body

		if resp.StatusCode != fiber.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			s.T().Errorf("Register Failed! Status: %d, Body: %s", resp.StatusCode, string(respBody))
		}

		s.Equal(fiber.StatusCreated, resp.StatusCode)
	})

	s.Run("Login_Success", func() {
		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		body, err := json.Marshal(payload)
		require.NoError(s.T(), err)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, 30000)
		require.NoError(s.T(), err)
		defer func() { _ = resp.Body.Close() }() // Bullet-proof: Tutup body

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[map[string]interface{}]

		// FIX: Tangkap error dari Decode agar lolos errcheck
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err, "Gagal mendecode JSON response login")

		s.True(result.Success)
		s.NotEmpty(result.Data["access_token"])
	})
}

func TestAuthIntegration(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
