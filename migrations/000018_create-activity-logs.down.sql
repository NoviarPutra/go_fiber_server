-- 1. Hapus Indeks
DROP INDEX IF EXISTS "idx_activity_logs_meta_gin";
DROP INDEX IF EXISTS "idx_activity_logs_entity_history";
DROP INDEX IF EXISTS "idx_activity_logs_user_feed";
DROP INDEX IF EXISTS "idx_activity_logs_company_feed";

-- 2. Hapus Tabel
DROP TABLE IF EXISTS "activity_logs";