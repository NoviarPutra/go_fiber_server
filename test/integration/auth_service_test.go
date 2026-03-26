package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type AuthServiceTestSuite struct {
	suite.Suite
	service *services.AuthService
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.Require().NotNil(testDBPool, "DB Pool harus terinisialisasi dari TestMain")
	s.service = services.NewAuthService(testDBPool)
}

// TearDownTest memastikan database bersih setiap kali menjalankan satu fungsi test
func (s *AuthServiceTestSuite) TearDownTest() {
	ctx := context.Background()
	testDBPool.Exec(ctx, "TRUNCATE users, refresh_tokens CASCADE")
}

// ─── REGISTER TESTS ──────────────────────────────────────────────────────────

func (s *AuthServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	req := &types.RegisterRequest{
		Email:    "new@officecore.id",
		Username: "newuser",
		Password: "Password123!",
	}

	res, err := s.service.Register(ctx, req)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Email, res.Email)
	s.NotEmpty(res.ID)
}

func (s *AuthServiceTestSuite) TestRegister_Duplicate() {
	ctx := context.Background()
	req := &types.RegisterRequest{
		Email:    "dup@officecore.id",
		Username: "dupuser",
		Password: "Password123!",
	}

	// Register pertama kali
	_, _ = s.service.Register(ctx, req)

	s.Run("Duplicate_Email", func() {
		req2 := *req
		req2.Username = "different_user"
		_, err := s.service.Register(ctx, &req2)
		s.ErrorIs(err, services.ErrEmailAlreadyExists)
	})

	s.Run("Duplicate_Username", func() {
		req3 := *req
		req3.Email = "different@email.id"
		_, err := s.service.Register(ctx, &req3)
		s.ErrorIs(err, services.ErrUsernameAlreadyExists)
	})
}

// ─── LOGIN TESTS ─────────────────────────────────────────────────────────────

func (s *AuthServiceTestSuite) TestLogin_Flow() {
	ctx := context.Background()
	email := "login@officecore.id"
	pass := "Secret123!"

	// Persiapan: Buat user dulu
	_, err := s.service.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "loginuser",
		Password: pass,
	})
	s.Require().NoError(err)

	s.Run("Success_Login", func() {
		req := &types.LoginRequest{Email: email, Password: pass}
		res, err := s.service.Login(ctx, req)

		s.NoError(err)
		s.NotEmpty(res.AccessToken)
		s.NotEmpty(res.RefreshToken)
	})

	s.Run("Wrong_Password", func() {
		req := &types.LoginRequest{Email: email, Password: "WrongPassword"}
		_, err := s.service.Login(ctx, req)
		s.ErrorIs(err, services.ErrInvalidCredentials)
	})

	s.Run("Non_Existent_User", func() {
		req := &types.LoginRequest{Email: "ghost@officecore.id", Password: pass}
		_, err := s.service.Login(ctx, req)
		s.ErrorIs(err, services.ErrInvalidCredentials)
	})
}

func (s *AuthServiceTestSuite) TestErrorConstants_Integrity() {
	ctx := context.Background()

	s.Run("Should_Return_ErrEmailAlreadyExists", func() {
		// 1. Setup: Register user pertama
		req := &types.RegisterRequest{
			Email:    "unique@officecore.id",
			Username: "uniqueuser",
			Password: "Password123!",
		}
		_, _ = s.service.Register(ctx, req)

		// 2. Action: Register dengan email yang sama
		_, err := s.service.Register(ctx, req)

		// 3. Assert: Pastikan error-nya IDENTIK dengan variable di services
		s.ErrorIs(err, services.ErrEmailAlreadyExists)
		s.Equal("email sudah terdaftar", err.Error())
	})

	s.Run("Should_Return_ErrInvalidCredentials_On_Login_Fail", func() {
		req := &types.LoginRequest{
			Email:    "wrong@officecore.id",
			Password: "AnyPassword",
		}
		_, err := s.service.Login(ctx, req)

		s.ErrorIs(err, services.ErrInvalidCredentials)
	})
}

func TestAuthServiceIntegration(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
