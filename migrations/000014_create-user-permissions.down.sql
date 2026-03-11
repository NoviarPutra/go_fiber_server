-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_user_permissions_modtime ON "user_permissions";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_user_permissions_lookup";
DROP INDEX IF EXISTS "idx_user_permission_unique";

-- 3. Hapus Tabel
DROP TABLE IF EXISTS "user_permissions";