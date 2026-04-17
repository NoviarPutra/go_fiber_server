package integration

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	app    *fiber.App
	secret string
}

func (s *AuthMiddlewareTestSuite) SetupSuite() {
	s.secret = "test-secret-key-123"
	os.Setenv("JWT_SECRET", s.secret)
	// Pastikan utils.AccessTokenExpiry terdefinisi, atau gunakan nilai manual
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.app = fiber.New()

	// Route dummy untuk proteksi
	s.app.Get("/protected", middlewares.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})
}

// Helper untuk membuat token manual untuk keperluan testing
func (s *AuthMiddlewareTestSuite) createTestToken(claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	s.Require().NoError(err)
	return t
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *AuthMiddlewareTestSuite) TestProtected_Success() {
	// 1. Generate token valid
	token := s.createTestToken(jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}, s.secret)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(200, resp.StatusCode)
}

func (s *AuthMiddlewareTestSuite) TestProtected_Cookie_Success() {
	token := s.createTestToken(jwt.MapClaims{
		"sub": "user-456",
		"exp": time.Now().Add(time.Hour).Unix(),
	}, s.secret)

	req := httptest.NewRequest("GET", "/protected", nil)
	// access_token is the internal name used in middlewares.Protected()
	req.Header.Set("Cookie", "access_token="+token)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(200, resp.StatusCode)
}

func (s *AuthMiddlewareTestSuite) TestProtected_MissingToken() {
	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(401, resp.StatusCode)
	// Pastikan pesan error sesuai dengan jwt_error_handler
}

func (s *AuthMiddlewareTestSuite) TestProtected_ExpiredToken() {
	// Token yang sudah mati 1 jam lalu
	token := s.createTestToken(jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}, s.secret)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := s.app.Test(req)
	s.Equal(401, resp.StatusCode)
}

func (s *AuthMiddlewareTestSuite) TestProtected_InvalidSignature() {
	// Token dibuat dengan secret yang BERBEDA
	token := s.createTestToken(jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}, "wrong-secret-key")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := s.app.Test(req)
	s.Equal(401, resp.StatusCode)
}

func (s *AuthMiddlewareTestSuite) TestProtected_MalformedToken() {
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ini-bukan-token-jwt-yang-benar")

	resp, _ := s.app.Test(req)
	s.Equal(401, resp.StatusCode)
}

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
