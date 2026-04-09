package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

type CompaniesService struct {
	db *pgxpool.Pool
}

func NewCompaniesService(db *pgxpool.Pool) *CompaniesService {
	return &CompaniesService{db: db}
}

func (s *CompaniesService) Create(ctx context.Context, req types.CreateCompanyRequest) (*types.CompanyRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO companies (name, code, logo_url)
		VALUES ($1, $2, $3)
		RETURNING id::text, name, code, logo_url, created_at, updated_at`

	var company types.CompanyRow
	err := utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, req.Name, req.Code, req.LogoUrl).Scan(
			&company.ID, &company.Name, &company.Code, &company.LogoUrl, &company.CreatedAt, &company.UpdatedAt,
		)
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrCompanyCodeExists
		}
		return nil, fmt.Errorf("gagal membuat data perusahaan")
	}

	return &company, nil
}

func (s *CompaniesService) GetAll(ctx context.Context, page, perPage int) ([]types.CompanyRow, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * perPage

	query := `
		WITH count_total AS (
			SELECT COUNT(*) as total FROM companies WHERE deleted_at IS NULL
		)
		SELECT 
			c.id::text, c.name, c.code, c.logo_url, 
			c.created_at, c.updated_at,
			ct.total
		FROM companies c, count_total ct
		WHERE c.deleted_at IS NULL
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data perusahaan")
	}
	defer rows.Close()

	var companies []types.CompanyRow
	var totalCount int64

	for rows.Next() {
		var c types.CompanyRow
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Code, &c.LogoUrl, &c.CreatedAt, &c.UpdatedAt, &totalCount,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, c)
	}

	if len(companies) == 0 {
		err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM companies WHERE deleted_at IS NULL").Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		return []types.CompanyRow{}, totalCount, nil
	}

	return companies, totalCount, nil
}

func (s *CompaniesService) GetByID(ctx context.Context, id string) (*types.CompanyRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id::text, name, code, logo_url, created_at, updated_at
		FROM companies
		WHERE id = $1 AND deleted_at IS NULL`

	var company types.CompanyRow
	err := s.db.QueryRow(ctx, query, id).Scan(
		&company.ID, &company.Name, &company.Code, &company.LogoUrl, &company.CreatedAt, &company.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("gagal mencari data perusahaan")
	}

	return &company, nil
}

func (s *CompaniesService) Update(ctx context.Context, id string, req types.UpdateCompanyRequest) (*types.CompanyRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Cek apakah eksis
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update data if not nil
	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}
	code := existing.Code
	if req.Code != nil {
		code = *req.Code
	}
	logoUrl := existing.LogoUrl
	if req.LogoUrl != nil {
		logoUrl = req.LogoUrl
	}

	query := `
		UPDATE companies
		SET name = $1, code = $2, logo_url = $3, updated_at = now()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id::text, name, code, logo_url, created_at, updated_at`

	var company types.CompanyRow
	err = utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, name, code, logoUrl, id).Scan(
			&company.ID, &company.Name, &company.Code, &company.LogoUrl, &company.CreatedAt, &company.UpdatedAt,
		)
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrCompanyCodeExists
		}
		return nil, fmt.Errorf("gagal merubah data perusahaan")
	}

	return &company, nil
}

func (s *CompaniesService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE companies SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
	var commandTag pgconn.CommandTag
	err := utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		var txErr error
		commandTag, txErr = tx.Exec(ctx, query, id)
		return txErr
	})
	if err != nil {
		return fmt.Errorf("gagal menghapus data perusahaan")
	}

	if commandTag.RowsAffected() == 0 {
		return ErrCompanyNotFound
	}

	return nil
}
