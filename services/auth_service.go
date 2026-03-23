package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/types"
	"github.com/yourusername/go_server/utils"
)

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{db: db}
}

// ─── Register ─────────────────────────────────────────────────────────────────

func (s *AuthService) Register(ctx context.Context, req *types.RegisterRequest) (*types.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.check_availability(ctx, req.Email, req.Username); err != nil {
		return nil, err
	}

	password_hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses password")
	}

	return s.create_user(ctx, req.Email, req.Username, password_hash)
}

func (s *AuthService) check_availability(ctx context.Context, email, username string) error {
	var email_exists, username_exists bool
	err := s.db.QueryRow(ctx,
		`SELECT
			EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL),
			EXISTS(SELECT 1 FROM users WHERE username = $2 AND deleted_at IS NULL)
		`, email, username,
	).Scan(&email_exists, &username_exists)
	if err != nil {
		return fmt.Errorf("gagal memeriksa ketersediaan akun")
	}
	if email_exists {
		return ErrEmailAlreadyExists
	}
	if username_exists {
		return ErrUsernameAlreadyExists
	}
	return nil
}

func (s *AuthService) create_user(ctx context.Context, email, username, password_hash string) (*types.RegisterResponse, error) {
	var user types.RegisterResponse
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (email, username, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id::text, email, username`,
		email, username, password_hash,
	).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat akun")
	}
	return &user, nil
}

// ─── Login ────────────────────────────────────────────────────────────────────

func (s *AuthService) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 1. Cari user
	user, err := s.find_user_by_email(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// 2. Verifikasi password
	match, err := utils.CheckPasswordHash(req.Password, user.password_hash)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	// 3. Generate tokens
	access_token, err := utils.GenerateAccessToken(user.id, user.email)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token")
	}

	refresh_token, expires_at, err := utils.GenerateRefreshToken(user.id)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token")
	}

	// 4. Simpan refresh token
	if err := s.save_refresh_token(ctx, user.id, refresh_token, expires_at); err != nil {
		return nil, err
	}

	// 5. Update last_login_at — tidak pakai goroutine, sudah ada timeout sendiri
	s.update_last_login(ctx, user.id)

	return &types.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()), // ✅ pakai konstanta exported
	}, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

type user_row struct {
	id            string
	email         string
	password_hash string
	is_active     bool
}

func (s *AuthService) find_user_by_email(ctx context.Context, email string) (*user_row, error) {
	var u user_row
	err := s.db.QueryRow(ctx,
		`SELECT id::text, email, password_hash, is_active
		 FROM users
		 WHERE email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(&u.id, &u.email, &u.password_hash, &u.is_active)
	if err != nil {
		return nil, ErrInvalidCredentials // generalisir — cegah user enumeration
	}
	if !u.is_active {
		return nil, ErrAccountInactive
	}
	return &u, nil
}

func (s *AuthService) save_refresh_token(ctx context.Context, user_id, token string, expires_at time.Time) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at)
		 VALUES ($1, $2, $3)`,
		user_id, token, expires_at,
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan sesi login")
	}
	return nil
}

func (s *AuthService) update_last_login(ctx context.Context, user_id string) {
	// Gagal update tidak batalkan login — intentional fire and forget
	s.db.Exec(ctx,
		`UPDATE users SET last_login_at = NOW() WHERE id = $1`,
		user_id,
	)
}
