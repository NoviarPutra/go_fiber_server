-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_teams_modtime ON "teams";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_teams_deleted_at";
DROP INDEX IF EXISTS "idx_teams_division_active";
DROP INDEX IF EXISTS "idx_teams_company_active";
DROP INDEX IF EXISTS "idx_teams_name_scoped_active_unique";

-- 3. Hapus Foreign Keys
ALTER TABLE "teams" DROP CONSTRAINT IF EXISTS "fk_teams_division";
ALTER TABLE "teams" DROP CONSTRAINT IF EXISTS "fk_teams_company";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "teams";