-- 1. Hapus Indeks
DROP INDEX IF EXISTS "idx_user_roles_role_id";
DROP INDEX IF EXISTS "idx_user_roles_company_user_id";
DROP INDEX IF EXISTS "idx_user_role_unique";

-- 2. Hapus Tabel (Foreign Keys otomatis terhapus)
DROP TABLE IF EXISTS "user_roles";