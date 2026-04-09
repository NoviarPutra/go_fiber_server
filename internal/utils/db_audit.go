package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditInfo digunakan untuk menyimpan informasi siapa yang melakukan mutasi HTTP
type AuditInfo struct {
	UserID    string
	CompanyID string
	IPAddress string
	UserAgent string
}

type contextKey string

const auditInfoKey contextKey = "audit_info"

// ContextWithAuditInfo dipakai untuk menyuntik manual (biasanya pada Testing atau RPC)
func ContextWithAuditInfo(ctx context.Context, info AuditInfo) context.Context {
	return context.WithValue(ctx, auditInfoKey, info)
}

// InjectAuditInfo digunakan di layer Handler/Middleware untuk menyimpan informasi request
func InjectAuditInfo(c *fiber.Ctx, userID, companyID string) context.Context {
	info := AuditInfo{
		UserID:    userID,
		CompanyID: companyID,
		IPAddress: c.IP(),
		UserAgent: string(c.Request().Header.UserAgent()),
	}
	return ContextWithAuditInfo(c.Context(), info)
}

// GetAuditInfo mengekstrak audit info dari context di layer Service
func GetAuditInfo(ctx context.Context) AuditInfo {
	if info, ok := ctx.Value(auditInfoKey).(AuditInfo); ok {
		return info
	}
	return AuditInfo{}
}

// WithAuditTx menjalankan transaksi database PostgreSQL di mana variabel-variabel
// context diset terlebih dahulu. Sehingga trigger database otomatis mencatat user_id & ip.
func WithAuditTx(ctx context.Context, pool *pgxpool.Pool, info AuditInfo, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err // Gagal memulai transaksi
	}
	// Pastikan rollback jika terjadi panic atau error tanpa commit
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Set local config variables (dijamin aman SQL Injection karena via arg queries)
	// false pada arg ke-3 artinya set_config TIDAK hanya read-only, dan
	// true diarg ke-4 (bisa disesuaikan di pgx jika param booleannya is_local) menunjuk "is_local=true" -> reset stlah trx beres.
	_, err = tx.Exec(ctx, `
		SELECT 
			set_config('app.audit_user_id', $1, true),
			set_config('app.audit_company_id', $2, true),
			set_config('app.audit_ip_address', $3, true),
			set_config('app.audit_user_agent', $4, true)
	`, info.UserID, info.CompanyID, info.IPAddress, info.UserAgent)

	if err != nil {
		return err // Gagal set local info
	}

	// Jalankan logika transaksi mutasi / service aslinya
	if err := fn(tx); err != nil {
		return err // fn() balikin error, bakal kena defer tx.Rollback() otomatis
	}

	// Commit data
	return tx.Commit(ctx)
}
