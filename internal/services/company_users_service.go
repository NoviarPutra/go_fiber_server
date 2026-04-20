package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

type CompanyUsersService struct {
	db *pgxpool.Pool
}

func NewCompanyUsersService(db *pgxpool.Pool) *CompanyUsersService {
	return &CompanyUsersService{db: db}
}

func (s *CompanyUsersService) AddUser(ctx context.Context, companyID uuid.UUID, req types.CreateCompanyUserRequest) (*types.CompanyUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO company_users (company_id, user_id, branch_id, employee_code, is_owner)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, company_id, user_id, branch_id, employee_code, is_owner, is_active, joined_at, updated_at`

	var cu types.CompanyUser
	err := utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, companyID, req.UserID, req.BranchID, req.EmployeeCode, req.IsOwner).Scan(
			&cu.ID, &cu.CompanyID, &cu.UserID, &cu.BranchID, &cu.EmployeeCode, &cu.IsOwner, &cu.IsActive, &cu.JoinedAt, &cu.UpdatedAt,
		)
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrCompanyUserAlreadyExists
		}
		return nil, fmt.Errorf("gagal menambahkan user ke perusahaan")
	}

	return &cu, nil
}

func (s *CompanyUsersService) List(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]types.CompanyUser, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * perPage

	query := `
		WITH count_total AS (
			SELECT COUNT(*) as total FROM company_users WHERE company_id = $1 AND deleted_at IS NULL
		)
		SELECT 
			cu.id, cu.company_id, cu.user_id, cu.branch_id, cu.employee_code, cu.is_owner, cu.is_active, 
			cu.joined_at, cu.updated_at, u.full_name, u.email,
			ct.total
		FROM company_users cu
		JOIN users u ON cu.user_id = u.id
		CROSS JOIN count_total ct
		WHERE cu.company_id = $1 AND cu.deleted_at IS NULL
		ORDER BY cu.joined_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(ctx, query, companyID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data user perusahaan")
	}
	defer rows.Close()

	var users []types.CompanyUser
	var totalCount int64

	for rows.Next() {
		var cu types.CompanyUser
		if err := rows.Scan(
			&cu.ID, &cu.CompanyID, &cu.UserID, &cu.BranchID, &cu.EmployeeCode, &cu.IsOwner, &cu.IsActive,
			&cu.JoinedAt, &cu.UpdatedAt, &cu.UserName, &cu.UserEmail, &totalCount,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, cu)
	}

	if len(users) == 0 {
		_ = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM company_users WHERE company_id = $1 AND deleted_at IS NULL", companyID).Scan(&totalCount)
		return []types.CompanyUser{}, totalCount, nil
	}

	return users, totalCount, nil
}

func (s *CompanyUsersService) GetDetail(ctx context.Context, companyID, userID uuid.UUID) (*types.CompanyUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			cu.id, cu.company_id, cu.user_id, cu.branch_id, cu.employee_code, cu.is_owner, cu.is_active, 
			cu.joined_at, cu.updated_at, u.full_name, u.email
		FROM company_users cu
		JOIN users u ON cu.user_id = u.id
		WHERE cu.company_id = $1 AND cu.user_id = $2 AND cu.deleted_at IS NULL`

	var cu types.CompanyUser
	err := s.db.QueryRow(ctx, query, companyID, userID).Scan(
		&cu.ID, &cu.CompanyID, &cu.UserID, &cu.BranchID, &cu.EmployeeCode, &cu.IsOwner, &cu.IsActive,
		&cu.JoinedAt, &cu.UpdatedAt, &cu.UserName, &cu.UserEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyUserNotFound
		}
		return nil, fmt.Errorf("gagal mengambil detail user perusahaan")
	}

	return &cu, nil
}

func (s *CompanyUsersService) Update(ctx context.Context, companyID, userID uuid.UUID, req types.UpdateCompanyUserRequest) (*types.CompanyUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get current
	existing, err := s.GetDetail(ctx, companyID, userID)
	if err != nil {
		return nil, err
	}

	// Prepare fields
	branchID := existing.BranchID
	if req.BranchID != nil {
		branchID = req.BranchID
	}
	employeeCode := existing.EmployeeCode
	if req.EmployeeCode != nil {
		employeeCode = req.EmployeeCode
	}
	isOwner := existing.IsOwner
	if req.IsOwner != nil {
		isOwner = *req.IsOwner
	}
	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	query := `
		UPDATE company_users
		SET branch_id = $1, employee_code = $2, is_owner = $3, is_active = $4, updated_at = now()
		WHERE company_id = $5 AND user_id = $6 AND deleted_at IS NULL
		RETURNING id, company_id, user_id, branch_id, employee_code, is_owner, is_active, joined_at, updated_at`

	var cu types.CompanyUser
	err = utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, branchID, employeeCode, isOwner, isActive, companyID, userID).Scan(
			&cu.ID, &cu.CompanyID, &cu.UserID, &cu.BranchID, &cu.EmployeeCode, &cu.IsOwner, &cu.IsActive, &cu.JoinedAt, &cu.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("gagal merubah data user perusahaan")
	}

	// Re-fetch to get user details
	return s.GetDetail(ctx, companyID, userID)
}

func (s *CompanyUsersService) Remove(ctx context.Context, companyID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE company_users SET deleted_at = now(), left_at = now() WHERE company_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var commandTag pgconn.CommandTag
	err := utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		var txErr error
		commandTag, txErr = tx.Exec(ctx, query, companyID, userID)
		return txErr
	})

	if err != nil {
		return fmt.Errorf("gagal menghapus user dari perusahaan")
	}

	if commandTag.RowsAffected() == 0 {
		return ErrCompanyUserNotFound
	}

	return nil
}
