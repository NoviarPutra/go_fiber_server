package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
)

type AuditLogsService struct {
	db *pgxpool.Pool
}

func NewAuditLogsService(db *pgxpool.Pool) *AuditLogsService {
	return &AuditLogsService{db: db}
}

func (s *AuditLogsService) GetAll(ctx context.Context, q types.AuditLogQuery, page, perPage int) ([]types.AuditLog, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	offset := (page - 1) * perPage

	// Base Query
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCounter := 1

	if q.CompanyID != "" {
		whereClause += fmt.Sprintf(" AND company_id = $%d", argCounter)
		args = append(args, q.CompanyID)
		argCounter++
	}
	if q.UserID != "" {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argCounter)
		args = append(args, q.UserID)
		argCounter++
	}
	if q.TableName != "" {
		whereClause += fmt.Sprintf(" AND table_name = $%d", argCounter)
		args = append(args, q.TableName)
		argCounter++
	}
	if q.RecordID != "" {
		whereClause += fmt.Sprintf(" AND record_id = $%d", argCounter)
		args = append(args, q.RecordID)
		argCounter++
	}

	// Hitung total dulu
	var totalCount int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data audit log")
	}

	if totalCount == 0 {
		return []types.AuditLog{}, 0, nil
	}

	// Fetch data
	query := fmt.Sprintf(`
		SELECT id::text, company_id::text, user_id::text, action, table_name, record_id::text,
		       old_data, new_data, ip_address, user_agent, created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argCounter, argCounter+1)

	args = append(args, perPage, offset)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data audit log")
	}
	defer rows.Close()

	var logs []types.AuditLog
	for rows.Next() {
		var log types.AuditLog
		if err := rows.Scan(
			&log.ID, &log.CompanyID, &log.UserID, &log.Action, &log.TableName, &log.RecordID,
			&log.OldData, &log.NewData, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, totalCount, nil
}

func (s *AuditLogsService) GetByID(ctx context.Context, id string) (*types.AuditLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id::text, company_id::text, user_id::text, action, table_name, record_id::text,
		       old_data, new_data, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE id = $1`

	var log types.AuditLog
	err := s.db.QueryRow(ctx, query, id).Scan(
		&log.ID, &log.CompanyID, &log.UserID, &log.Action, &log.TableName, &log.RecordID,
		&log.OldData, &log.NewData, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrAuditLogNotFound
		}
		return nil, fmt.Errorf("gagal mencari data audit log")
	}

	return &log, nil
}
