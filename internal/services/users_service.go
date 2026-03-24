package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
)

type UsersService struct {
	db *pgxpool.Pool
}

func NewUsersService(db *pgxpool.Pool) *UsersService {
	return &UsersService{db: db}
}

func (s *UsersService) GetUsers(ctx context.Context, page, per_page int) ([]types.UserRow, int64, error) {
	// Timeout guard
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * per_page

	// 1. Ambil total count & data dalam 1 query pakai window function
	// Lebih efisien dari 2 query terpisah (SELECT COUNT + SELECT data)
	rows, err := s.db.Query(ctx,
		`SELECT
			id::text,
			email,
			username,
			is_active,
			created_at,
			last_login_at,
			COUNT(*) OVER() AS total_count
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		per_page, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data users")
	}
	defer rows.Close()

	// 2. Scan hasil
	var users []types.UserRow
	var total_count int64

	for rows.Next() {
		var u types.UserRow
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.IsActive,
			&u.CreatedAt,
			&u.LastLoginAt,
			&total_count,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memproses data users")
		}
		users = append(users, u)
	}

	// 3. Cek error setelah iterasi
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("gagal membaca data users")
	}

	// 4. Jika tidak ada data, kembalikan slice kosong (bukan nil)
	if users == nil {
		users = []types.UserRow{}
	}

	return users, total_count, nil
}
