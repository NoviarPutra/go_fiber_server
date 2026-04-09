package types

import (
	"encoding/json"
	"time"
)

// AuditLog merepresentasikan struktur data dari tabel audit_logs.
type AuditLog struct {
	ID        string          `json:"id"`
	CompanyID *string         `json:"company_id,omitempty"`
	UserID    *string         `json:"user_id,omitempty"`
	Action    string          `json:"action"`
	TableName string          `json:"table_name"`
	RecordID  *string         `json:"record_id,omitempty"`
	OldData   json.RawMessage `json:"old_data,omitempty"`
	NewData   json.RawMessage `json:"new_data,omitempty"`
	IPAddress *string         `json:"ip_address,omitempty"`
	UserAgent *string         `json:"user_agent,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// AuditLogQuery digunakan untuk filter query parameter di GET /api/v1/audit-logs
type AuditLogQuery struct {
	CompanyID string `query:"company_id" validate:"omitempty,uuid"`
	UserID    string `query:"user_id" validate:"omitempty,uuid"`
	TableName string `query:"table_name" validate:"omitempty"`
	RecordID  string `query:"record_id" validate:"omitempty,uuid"`
}
