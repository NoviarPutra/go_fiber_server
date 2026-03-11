-- 1. Hapus Indeks
DROP INDEX IF EXISTS "idx_role_permissions_allowed";
DROP INDEX IF EXISTS "idx_role_permissions_permission_id";
DROP INDEX IF EXISTS "idx_role_permission_unique";

-- 2. Hapus Tabel (Foreign Keys otomatis terhapus)
DROP TABLE IF EXISTS "role_permissions";