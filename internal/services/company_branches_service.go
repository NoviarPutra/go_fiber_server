package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
)

type CompanyBranchesService struct {
	db *pgxpool.Pool
}

func NewCompanyBranchesService(db *pgxpool.Pool) *CompanyBranchesService {
	return &CompanyBranchesService{db: db}
}

func (s *CompanyBranchesService) Create(ctx context.Context, req types.CreateCompanyBranchRequest) (*types.CompanyBranchRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO company_branches (company_id, name, address, timezone)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, company_id::text, name, address, timezone, created_at, updated_at`

	var branch types.CompanyBranchRow
	err := s.db.QueryRow(ctx, query, req.CompanyID, req.Name, req.Address, req.Timezone).Scan(
		&branch.ID, &branch.CompanyID, &branch.Name, &branch.Address, &branch.Timezone, &branch.CreatedAt, &branch.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrCompanyBranchNameExists
		}
		return nil, fmt.Errorf("gagal membuat data cabang perusahaan: %w", err)
	}

	return &branch, nil
}

func (s *CompanyBranchesService) GetAll(ctx context.Context, page, perPage int) ([]types.CompanyBranchRow, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * perPage

	query := `
		WITH count_total AS (
			SELECT COUNT(*) as total FROM company_branches WHERE deleted_at IS NULL
		)
		SELECT 
			cb.id::text, cb.company_id::text, cb.name, cb.address, cb.timezone, 
			cb.created_at, cb.updated_at,
			ct.total
		FROM company_branches cb, count_total ct
		WHERE cb.deleted_at IS NULL
		ORDER BY cb.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data cabang perusahaan: %w", err)
	}
	defer rows.Close()

	var branches []types.CompanyBranchRow
	var totalCount int64

	for rows.Next() {
		var cb types.CompanyBranchRow
		if err := rows.Scan(
			&cb.ID, &cb.CompanyID, &cb.Name, &cb.Address, &cb.Timezone, &cb.CreatedAt, &cb.UpdatedAt, &totalCount,
		); err != nil {
			return nil, 0, err
		}
		branches = append(branches, cb)
	}

	if len(branches) == 0 {
		err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM company_branches WHERE deleted_at IS NULL").Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		return []types.CompanyBranchRow{}, totalCount, nil
	}

	return branches, totalCount, nil
}

func (s *CompanyBranchesService) GetByID(ctx context.Context, id string) (*types.CompanyBranchRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id::text, company_id::text, name, address, timezone, created_at, updated_at
		FROM company_branches
		WHERE id = $1 AND deleted_at IS NULL`

	var branch types.CompanyBranchRow
	err := s.db.QueryRow(ctx, query, id).Scan(
		&branch.ID, &branch.CompanyID, &branch.Name, &branch.Address, &branch.Timezone, &branch.CreatedAt, &branch.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyBranchNotFound
		}
		return nil, fmt.Errorf("gagal mencari data cabang perusahaan: %w", err)
	}

	return &branch, nil
}

func (s *CompanyBranchesService) Update(ctx context.Context, id string, req types.UpdateCompanyBranchRequest) (*types.CompanyBranchRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Cek apakah eksis
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update data if not nil
	companyID := existing.CompanyID
	if req.CompanyID != nil {
		companyID = *req.CompanyID
	}
	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}
	address := existing.Address
	if req.Address != nil {
		address = req.Address
	}
	timezone := existing.Timezone
	if req.Timezone != nil {
		timezone = *req.Timezone
	}

	query := `
		UPDATE company_branches
		SET company_id = $1, name = $2, address = $3, timezone = $4, updated_at = now()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id::text, company_id::text, name, address, timezone, created_at, updated_at`

	var branch types.CompanyBranchRow
	err = s.db.QueryRow(ctx, query, companyID, name, address, timezone, id).Scan(
		&branch.ID, &branch.CompanyID, &branch.Name, &branch.Address, &branch.Timezone, &branch.CreatedAt, &branch.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrCompanyBranchNameExists
		}
		return nil, fmt.Errorf("gagal merubah data cabang perusahaan: %w", err)
	}

	return &branch, nil
}

func (s *CompanyBranchesService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE company_branches SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
	commandTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data cabang perusahaan: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrCompanyBranchNotFound
	}

	return nil
}
