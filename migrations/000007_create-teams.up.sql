CREATE TABLE IF NOT EXISTS "teams" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "division_id" uuid,
    "name" text NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,

    -- Constraint minimal untuk validasi nama
    CONSTRAINT "teams_name_check" CHECK (char_length(name) >= 1)
);

-- 1. Foreign Keys dengan Integritas Data
ALTER TABLE "teams" 
ADD CONSTRAINT "fk_teams_company" FOREIGN KEY ("company_id") REFERENCES "companies"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_teams_division" FOREIGN KEY ("division_id") REFERENCES "divisions"("id") ON DELETE SET NULL;

-- 2. Robust Unique Index (Opsional tetapi direkomendasikan)
-- Mencegah nama tim yang sama di dalam satu divisi/perusahaan (hanya untuk data aktif)
CREATE UNIQUE INDEX "idx_teams_name_scoped_active_unique" 
ON "teams" ("company_id", coalesce("division_id", '00000000-0000-0000-0000-000000000000'), "name") 
WHERE "deleted_at" IS NULL;

-- 3. Optimized Lookup Index untuk Perusahaan (Data Aktif)
CREATE INDEX "idx_teams_company_active" 
ON "teams" ("company_id") 
WHERE "deleted_at" IS NULL;

-- 4. Optimized Lookup Index untuk Divisi (Data Aktif)
CREATE INDEX "idx_teams_division_active" 
ON "teams" ("division_id") 
WHERE "division_id" IS NOT NULL AND "deleted_at" IS NULL;

-- 5. Optimized Index untuk Soft Delete
CREATE INDEX "idx_teams_deleted_at" 
ON "teams" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 6. Trigger updated_at Otomatis
CREATE TRIGGER update_teams_modtime
    BEFORE UPDATE ON "teams"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();