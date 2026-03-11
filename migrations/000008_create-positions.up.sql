CREATE TABLE IF NOT EXISTS "positions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "name" text NOT NULL,
    "level" integer DEFAULT 1 NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL, -- Tambahkan ini untuk konsistensi
    "deleted_at" timestamp with time zone,

    -- Constraint untuk validasi data
    CONSTRAINT "positions_name_check" CHECK (char_length(name) >= 1),
    CONSTRAINT "positions_level_check" CHECK (level >= 0)
);

-- 1. Foreign Key dengan Cascade Delete
ALTER TABLE "positions" 
ADD CONSTRAINT "fk_positions_company" 
FOREIGN KEY ("company_id") REFERENCES "companies"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Robust Unique Index untuk Nama Jabatan
-- Mencegah nama jabatan ganda di perusahaan yang sama (hanya untuk data aktif)
CREATE UNIQUE INDEX "idx_positions_name_company_active_unique" 
ON "positions" ("company_id", "name") 
WHERE "deleted_at" IS NULL;

-- 3. Optimized Index untuk Look-up Perusahaan & Sorting Level
-- Sangat berguna untuk query: SELECT * FROM positions WHERE company_id = ? ORDER BY level ASC
CREATE INDEX "idx_positions_company_level_active" 
ON "positions" ("company_id", "level") 
WHERE "deleted_at" IS NULL;

-- 4. Optimized Index untuk Soft Delete
CREATE INDEX "idx_positions_deleted_at" 
ON "positions" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 5. Trigger updated_at Otomatis
CREATE TRIGGER update_positions_modtime
    BEFORE UPDATE ON "positions"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();