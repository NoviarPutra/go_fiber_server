-- 1. Hapus Indeks
DROP INDEX IF EXISTS "idx_audit_logs_created_at_desc";
DROP INDEX IF EXISTS "idx_audit_logs_record_lookup";
DROP INDEX IF EXISTS "idx_audit_logs_user_lookup";
DROP INDEX IF EXISTS "idx_audit_logs_company_lookup";

-- 2. Hapus Tabel
DROP TABLE IF EXISTS "audit_logs";