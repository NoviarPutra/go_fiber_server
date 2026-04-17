package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/auth"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

type TokenManagementTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *TokenManagementTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.app.Use(middlewares.DBMiddleware(testDBPool))
	
	auth_group := s.app.Group("/api/v1/auth")
	auth_group.Post("/login", auth.Login)
	auth_group.Post("/refresh", auth.Refresh)
	auth_group.Post("/revoke", auth.Revoke)
}

func (s *TokenManagementTestSuite) TearDownTest() {
	_, _ = testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
}

// ─── HELPERS ──────────────────────────────────────────────────────────────────

func (s *TokenManagementTestSuite) seedUser(email, password string, active bool) string {
	hash, err := utils.HashPassword(password)
	s.Require().NoError(err)

	var userID string
	query := `INSERT INTO users (email, username, password_hash, is_active) 
              VALUES ($1, $1, $2, $3) RETURNING id`
	err = testDBPool.QueryRow(context.Background(), query, email, hash, active).Scan(&userID)
	s.Require().NoError(err)
	return userID
}

func (s *TokenManagementTestSuite) login(email, password string) *types.LoginResponse {
	payload := types.LoginRequest{Email: email, Password: password}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := s.app.Test(req, 10000)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	
	data, ok := result["data"].(map[string]interface{})
	s.Require().True(ok, "Format data response tidak sesuai: %v", result)

	return &types.LoginResponse{
		AccessToken:  data["access_token"].(string),
		RefreshToken: data["refresh_token"].(string),
		ExpiresIn:    int64(data["expires_in"].(float64)),
	}
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *TokenManagementTestSuite) TestRefresh_Success() {
	email := "test@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	payload := types.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Wait a bit to ensure timestamps differ if needed (though not strictly necessary here)
	time.Sleep(100 * time.Millisecond)
	
	resp, err := s.app.Test(req, 5000)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(200, resp.StatusCode)
	
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	s.Require().NoError(err)
	s.Equal("Token berhasil diperbarui", res["message"])
	
	data, ok := res["data"].(map[string]interface{})
	s.Require().True(ok, "Format data response tidak sesuai: %v", res)

	s.NotEmpty(data["access_token"])
	s.NotEmpty(data["refresh_token"])
	s.NotEqual(loginResp.RefreshToken, data["refresh_token"], "Refresh token should be rotated")
}

func (s *TokenManagementTestSuite) TestRefresh_RotateAndRevokeOld() {
	email := "rotation@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Refresh first time
	payload := types.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.Test(req, 5000)
	s.Equal(200, resp.StatusCode)
	resp.Body.Close()

	// Try refresh with the OLD token again (Reuse Attack)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := s.app.Test(req2, 5000)
	s.Equal(401, resp2.StatusCode, "Old token should be revoked after rotation")
	resp2.Body.Close()
}

func (s *TokenManagementTestSuite) TestRevoke_Success() {
	email := "revoke@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	payload := types.RevokeTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/api/v1/auth/revoke", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.app.Test(req, 5000)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(200, resp.StatusCode)

	// Verify it's revoked by trying to refresh
	refreshPayload := types.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	rBody, _ := json.Marshal(refreshPayload)
	rReq := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(rBody))
	rReq.Header.Set("Content-Type", "application/json")
	rResp, _ := s.app.Test(rReq, 5000)
	s.Equal(401, rResp.StatusCode)
	rResp.Body.Close()
}

func (s *TokenManagementTestSuite) TestRefresh_InvalidToken() {
	payload := types.RefreshTokenRequest{RefreshToken: "invalid-token-here"}
	body, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := s.app.Test(req, 5000)
	s.Equal(401, resp.StatusCode)
	resp.Body.Close()
}

func (s *TokenManagementTestSuite) TestRefresh_InactiveUser() {
	email := "inactive@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Deactivate user
	_, _ = testDBPool.Exec(context.Background(), "UPDATE users SET is_active = false WHERE email = $1", email)

	payload := types.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := s.app.Test(req, 5000)
	s.Equal(401, resp.StatusCode, "Should fail if user is no longer active")
	resp.Body.Close()
}

func (s *TokenManagementTestSuite) TestRevoke_AlreadyRevoked() {
	email := "already-revoked@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	payload := types.RevokeTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	
	// Revoke first time - Success
	req1 := httptest.NewRequest("POST", "/api/v1/auth/revoke", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := s.app.Test(req1, 5000)
	s.Equal(200, resp1.StatusCode)
	resp1.Body.Close()

	// Revoke second time - Should return 401 (ErrRefreshTokenInvalid)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/revoke", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := s.app.Test(req2, 5000)
	s.Equal(401, resp2.StatusCode)
	resp2.Body.Close()
}

func (s *TokenManagementTestSuite) TestRefresh_ExpiredToken() {
	email := "expired@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Manually expire the token in DB
	_, err := testDBPool.Exec(context.Background(), 
		"UPDATE refresh_tokens SET expires_at = NOW() - INTERVAL '1 hour' WHERE token = $1", 
		loginResp.RefreshToken)
	s.Require().NoError(err)

	payload := types.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	body, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := s.app.Test(req, 5000)
	s.Equal(401, resp.StatusCode, "Should fail if token is expired")
	resp.Body.Close()
}

func TestTokenManagement(t *testing.T) {
	suite.Run(t, new(TokenManagementTestSuite))
}
