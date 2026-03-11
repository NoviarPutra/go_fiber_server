-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_roles_modtime ON "roles";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_roles_deleted_at";
DROP INDEX IF EXISTS "idx_roles_system_lookup";
DROP INDEX IF EXISTS "idx_roles_company_active";
DROP INDEX IF EXISTS "idx_roles_name_company_active_unique";

-- 3. Hapus Foreign Key
ALTER TABLE "roles" DROP CONSTRAINT IF EXISTS "fk_roles_company";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "roles";