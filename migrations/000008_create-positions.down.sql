-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_positions_modtime ON "positions";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_positions_deleted_at";
DROP INDEX IF EXISTS "idx_positions_company_level_active";
DROP INDEX IF EXISTS "idx_positions_name_company_active_unique";

-- 3. Hapus Foreign Key
ALTER TABLE "positions" DROP CONSTRAINT IF EXISTS "fk_positions_company";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "positions";