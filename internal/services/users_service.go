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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * per_page

	// Gunakan CTE (Common Table Expression) untuk mendapatkan total & data sekaligus
	// Ini menjamin total_count tetap muncul meski data page kosong
	query := `
        WITH count_total AS (
            SELECT COUNT(*) as total FROM users WHERE deleted_at IS NULL
        )
        SELECT 
            u.id::text, u.email, u.username, u.is_active, 
            u.created_at, u.last_login_at,
            ct.total
        FROM users u, count_total ct
        WHERE u.deleted_at IS NULL
        ORDER BY u.created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, per_page, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data users")
	}
	defer rows.Close()

	var users []types.UserRow
	var total_count int64

	for rows.Next() {
		var u types.UserRow
		// Scan tetap sama
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.IsActive, &u.CreatedAt, &u.LastLoginAt, &total_count); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	// --- FIX CRITICAL DISINI ---
	// Jika users kosong karena offset terlalu besar,
	// kita ambil count secara terpisah hanya jika rows.Next() tidak pernah jalan
	if len(users) == 0 {
		err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&total_count)
		if err != nil {
			return nil, 0, err
		}
		return []types.UserRow{}, total_count, nil
	}

	return users, total_count, nil
}
