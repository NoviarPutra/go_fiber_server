-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_users_modtime ON "users";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_users_deleted_at";
DROP INDEX IF EXISTS "idx_users_auth_lookup";
DROP INDEX IF EXISTS "idx_users_username_active_unique";
DROP INDEX IF EXISTS "idx_users_email_active_unique";

-- 3. Hapus Tabel
DROP TABLE IF EXISTS "users";