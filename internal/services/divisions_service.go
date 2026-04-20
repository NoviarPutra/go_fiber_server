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

type DivisionsService struct {
	db *pgxpool.Pool
}

func NewDivisionsService(db *pgxpool.Pool) *DivisionsService {
	return &DivisionsService{db: db}
}

func (s *DivisionsService) Create(ctx context.Context, req types.CreateDivisionRequest) (*types.DivisionRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Verify company exists
	var companyID string
	err := s.db.QueryRow(ctx, "SELECT id::text FROM companies WHERE id = $1 AND deleted_at IS NULL", req.CompanyID).Scan(&companyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("gagal verifikasi data perusahaan")
	}

	query := `
		INSERT INTO divisions (company_id, name, code)
		VALUES ($1, $2, $3)
		RETURNING id::text, company_id::text, name, code, created_at, updated_at`

	var division types.DivisionRow
	err = utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, req.CompanyID, req.Name, req.Code).Scan(
			&division.ID, &division.CompanyID, &division.Name, &division.Code, &division.CreatedAt, &division.UpdatedAt,
		)
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrDivisionCodeExists
		}
		return nil, fmt.Errorf("gagal membuat data divisi")
	}

	return &division, nil
}

func (s *DivisionsService) GetAll(ctx context.Context, companyID string, page, perPage int) ([]types.DivisionRow, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * perPage

	query := `
		WITH count_total AS (
			SELECT COUNT(*) as total FROM divisions WHERE company_id = $1 AND deleted_at IS NULL
		)
		SELECT 
			d.id::text, d.company_id::text, d.name, d.code, 
			d.created_at, d.updated_at,
			ct.total
		FROM divisions d, count_total ct
		WHERE d.company_id = $1 AND d.deleted_at IS NULL
		ORDER BY d.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(ctx, query, companyID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data divisi")
	}
	defer rows.Close()

	var divisions []types.DivisionRow
	var totalCount int64

	for rows.Next() {
		var d types.DivisionRow
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt, &totalCount,
		); err != nil {
			return nil, 0, err
		}
		divisions = append(divisions, d)
	}

	if len(divisions) == 0 {
		err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM divisions WHERE company_id = $1 AND deleted_at IS NULL", companyID).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
		return []types.DivisionRow{}, totalCount, nil
	}

	return divisions, totalCount, nil
}

func (s *DivisionsService) GetByID(ctx context.Context, id string) (*types.DivisionRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id::text, company_id::text, name, code, created_at, updated_at
		FROM divisions
		WHERE id = $1 AND deleted_at IS NULL`

	var division types.DivisionRow
	err := s.db.QueryRow(ctx, query, id).Scan(
		&division.ID, &division.CompanyID, &division.Name, &division.Code, &division.CreatedAt, &division.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDivisionNotFound
		}
		return nil, fmt.Errorf("gagal mencari data divisi")
	}

	return &division, nil
}

func (s *DivisionsService) Update(ctx context.Context, id string, req types.UpdateDivisionRequest) (*types.DivisionRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}
	code := existing.Code
	if req.Code != nil {
		code = req.Code
	}

	query := `
		UPDATE divisions
		SET name = $1, code = $2, updated_at = now()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id::text, company_id::text, name, code, created_at, updated_at`

	var division types.DivisionRow
	err = utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, name, code, id).Scan(
			&division.ID, &division.CompanyID, &division.Name, &division.Code, &division.CreatedAt, &division.UpdatedAt,
		)
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrDivisionCodeExists
		}
		return nil, fmt.Errorf("gagal merubah data divisi")
	}

	return &division, nil
}

func (s *DivisionsService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE divisions SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
	var commandTag pgconn.CommandTag
	err := utils.WithAuditTx(ctx, s.db, utils.GetAuditInfo(ctx), func(tx pgx.Tx) error {
		var txErr error
		commandTag, txErr = tx.Exec(ctx, query, id)
		return txErr
	})
	if err != nil {
		return fmt.Errorf("gagal menghapus data divisi")
	}

	if commandTag.RowsAffected() == 0 {
		return ErrDivisionNotFound
	}

	return nil
}
