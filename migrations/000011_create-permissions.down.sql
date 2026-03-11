-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_permissions_modtime ON "permissions";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_permissions_code_lookup";
DROP INDEX IF EXISTS "idx_permissions_module";

-- 3. Hapus Tabel
DROP TABLE IF EXISTS "permissions";