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

	loginResp := &types.LoginResponse{
		ExpiresIn: int64(data["expires_in"].(float64)),
	}

	// Extract tokens from Cookies (HTTPOnly)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == utils.CookieAccessToken {
			loginResp.AccessToken = cookie.Value
		}
		if cookie.Name == utils.CookieRefreshToken {
			loginResp.RefreshToken = cookie.Value
		}
	}

	s.NotEmpty(loginResp.AccessToken, "Access token harus ada di cookie")
	s.NotEmpty(loginResp.RefreshToken, "Refresh token harus ada di cookie")

	return loginResp
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *TokenManagementTestSuite) TestRefresh_Success() {
	email := "test@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Use Cookie instead of body
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)

	// Wait a bit to ensure tokens differ (uuid/jti handles this now but good practice)
	time.Sleep(100 * time.Millisecond)

	resp, err := s.app.Test(req, 5000)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(200, resp.StatusCode)

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	s.Require().NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.Require().True(ok, "Format data response tidak sesuai: %v", res)

	s.Empty(data["access_token"], "Access token tidak boleh ada di body")
	s.Empty(data["refresh_token"], "Refresh token tidak boleh ada di body")

	// Extract from cookies
	var newAccessToken, newRefreshToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == utils.CookieAccessToken {
			newAccessToken = cookie.Value
		}
		if cookie.Name == utils.CookieRefreshToken {
			newRefreshToken = cookie.Value
		}
	}

	s.NotEmpty(newAccessToken, "Access token harus ada di cookie")
	s.NotEmpty(newRefreshToken, "Refresh token harus ada di cookie")
	s.NotEqual(loginResp.RefreshToken, newRefreshToken, "Refresh token harus di-rotate")

	// ----------------------------------------------------
	// Coverage untuk Mobile
	// ----------------------------------------------------
	reqMobile := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	reqMobile.Header.Set("Cookie", utils.CookieRefreshToken+"="+newRefreshToken) // pakai token baru
	reqMobile.Header.Set("X-Client-Type", "mobile")

	respMobile, _ := s.app.Test(reqMobile, 5000)
	s.Equal(200, respMobile.StatusCode)

	var resMobile map[string]interface{}
	err = json.NewDecoder(respMobile.Body).Decode(&resMobile)
	s.NoError(err)
	respMobile.Body.Close()

	dataMobile, ok := resMobile["data"].(map[string]interface{})
	s.True(ok)
	s.NotEmpty(dataMobile["access_token"], "Access Token Refresh harus dikembalikan pada response JSON untuk Mobile Client")
	s.NotEmpty(dataMobile["refresh_token"], "Refresh Token Refresh harus dikembalikan pada response JSON untuk Mobile Client")
}

func (s *TokenManagementTestSuite) TestRefresh_RotateAndRevokeOld() {
	email := "rotation@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Refresh first time with cookie
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)
	resp, _ := s.app.Test(req, 5000)
	s.Equal(200, resp.StatusCode)
	resp.Body.Close()

	// Try refresh with the OLD token again (Reuse Attack)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req2.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)
	resp2, _ := s.app.Test(req2, 5000)
	s.Equal(401, resp2.StatusCode, "Old token should be revoked after rotation")
	resp2.Body.Close()
}

func (s *TokenManagementTestSuite) TestRevoke_Success() {
	email := "revoke@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Revoke with cookie
	req := httptest.NewRequest("POST", "/api/v1/auth/revoke", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)

	resp, err := s.app.Test(req, 5000)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(200, resp.StatusCode)

	// Verify it's revoked by trying to refresh with cookie
	rReq := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	rReq.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)
	rResp, _ := s.app.Test(rReq, 5000)
	s.Equal(401, rResp.StatusCode)
	rResp.Body.Close()

	// Verify identity cookie is cleared
	var accessTokenCleared bool
	for _, c := range resp.Cookies() {
		if c.Name == utils.CookieAccessToken && c.Value == "" {
			accessTokenCleared = true
		}
	}
	s.True(accessTokenCleared, "Auth cookie should be cleared on logout")
}

func (s *TokenManagementTestSuite) TestRefresh_InvalidToken() {
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"=invalid-token")

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

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)

	resp, _ := s.app.Test(req, 5000)
	s.Equal(401, resp.StatusCode, "Should fail if user is no longer active")
	resp.Body.Close()
}

func (s *TokenManagementTestSuite) TestRevoke_AlreadyRevoked() {
	email := "already-revoked@example.com"
	pass := "Password123!"
	s.seedUser(email, pass, true)
	loginResp := s.login(email, pass)

	// Revoke via cookie
	req1 := httptest.NewRequest("POST", "/api/v1/auth/revoke", nil)
	req1.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)
	resp1, _ := s.app.Test(req1, 5000)
	s.Equal(200, resp1.StatusCode)
	resp1.Body.Close()

	// Revoke second time via cookie
	req2 := httptest.NewRequest("POST", "/api/v1/auth/revoke", nil)
	req2.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)
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

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Cookie", utils.CookieRefreshToken+"="+loginResp.RefreshToken)

	resp, _ := s.app.Test(req, 5000)
	s.Equal(401, resp.StatusCode, "Should fail if token is expired")
	resp.Body.Close()
}

func TestTokenManagement(t *testing.T) {
	suite.Run(t, new(TokenManagementTestSuite))
}
