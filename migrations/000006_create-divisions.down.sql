-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_divisions_modtime ON "divisions";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_divisions_deleted_at";
DROP INDEX IF EXISTS "idx_divisions_company_active";
DROP INDEX IF EXISTS "idx_divisions_code_company_active_unique";

-- 3. Hapus Foreign Key
ALTER TABLE "divisions" DROP CONSTRAINT IF EXISTS "fk_divisions_company";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "divisions";