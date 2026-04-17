package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
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
	s.UpdateLastLogin(ctx, user.id)

	return &types.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
	}, nil
}

// ─── Refresh Token ────────────────────────────────────────────────────────────

func (s *AuthService) Refresh(ctx context.Context, refresh_token string) (*types.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Validasi token dan ambil data user (di dalam TX)
	var userIDStr, oldTokenIDStr string
	var email string
	err = tx.QueryRow(ctx,
		`SELECT rt.id::text, rt.user_id::text, u.email
		 FROM refresh_tokens rt
		 JOIN users u ON u.id = rt.user_id
		 WHERE rt.token = $1 
		   AND rt.revoked_at IS NULL 
		   AND rt.expires_at > NOW()
		   AND u.deleted_at IS NULL 
		   AND u.is_active = true`,
		refresh_token,
	).Scan(&oldTokenIDStr, &userIDStr, &email)

	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// 2. Generate token baru
	new_access_token, err := utils.GenerateAccessToken(userIDStr, email)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	new_refresh_token, expires_at, err := utils.GenerateRefreshToken(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token: %w", err)
	}

	// 3. Simpan token baru
	newTokenID := uuid.New().String()
	_, err = tx.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		 VALUES ($1::uuid, $2::uuid, $3, $4)`,
		newTokenID, userIDStr, new_refresh_token, expires_at,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan refresh token baru: %w", err)
	}

	// 4. Revoke yang lama
	_, err = tx.Exec(ctx,
		`UPDATE refresh_tokens 
		 SET revoked_at = NOW(), replaced_by = $1::uuid
		 WHERE id = $2::uuid`,
		newTokenID, oldTokenIDStr,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mencabut refresh token lama: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit transaksi")
	}

	return &types.LoginResponse{
		AccessToken:  new_access_token,
		RefreshToken: new_refresh_token,
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
	}, nil
}

// ─── Revoke Token ─────────────────────────────────────────────────────────────

func (s *AuthService) Revoke(ctx context.Context, refresh_token string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.db.Exec(ctx,
		`UPDATE refresh_tokens 
		 SET revoked_at = NOW() 
		 WHERE token = $1 AND revoked_at IS NULL`,
		refresh_token,
	)
	if err != nil {
		return fmt.Errorf("gagal mencabut token")
	}

	if result.RowsAffected() == 0 {
		return ErrRefreshTokenInvalid
	}

	return nil
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

func (s *AuthService) save_refresh_token(ctx context.Context, user_id string, token string, expires_at time.Time) error {
	id := uuid.New()
	_, err := s.db.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		id, user_id, token, expires_at,
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan refresh token: %w", err)
	}
	return nil
}

func (s *AuthService) UpdateLastLogin(ctx context.Context, user_id string) {
	// Gunakan _ untuk memberitahu linter bahwa kita sadar ada error tapi memilih mengabaikannya
	// Namun, sangat disarankan untuk setidaknya log jika terjadi error
	_, err := s.db.Exec(ctx,
		`UPDATE users SET last_login_at = NOW() WHERE id = $1`,
		user_id,
	)

	if err != nil {
		// Log saja, jangan return error karena ini non-critical
		log.Printf("Warning: gagal update last_login untuk user %s: %v\n", user_id, err)
	}
}
