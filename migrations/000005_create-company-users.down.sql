-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_company_users_modtime ON "company_users";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_cu_deleted_at";
DROP INDEX IF EXISTS "idx_cu_branch_lookup";
DROP INDEX IF EXISTS "idx_cu_auth_lookup";
DROP INDEX IF EXISTS "idx_cu_company_user_active_unique";

-- 3. Hapus Tabel (Otomatis menghapus FK)
DROP TABLE IF EXISTS "company_users";