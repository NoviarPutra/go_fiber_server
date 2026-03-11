CREATE TABLE IF NOT EXISTS "divisions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "name" text NOT NULL,
    "code" text,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,

    -- Constraint untuk integritas data
    CONSTRAINT "divisions_name_check" CHECK (char_length(name) >= 1)
);

-- 1. Foreign Key dengan Cascade Delete
ALTER TABLE "divisions" 
ADD CONSTRAINT "fk_divisions_company" 
FOREIGN KEY ("company_id") REFERENCES "companies"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Robust Unique Index untuk Code per Company
-- Mencegah kode divisi ganda di perusahaan yang sama (hanya untuk data aktif)
CREATE UNIQUE INDEX "idx_divisions_code_company_active_unique" 
ON "divisions" ("company_id", "code") 
WHERE "deleted_at" IS NULL AND "code" IS NOT NULL;

-- 3. Optimized Index untuk Look-up Perusahaan (Data Aktif)
-- Sangat kencang untuk dropdown divisi di sisi Frontend
CREATE INDEX "idx_divisions_company_active" 
ON "divisions" ("company_id") 
WHERE "deleted_at" IS NULL;

-- 4. Optimized Index untuk Soft Delete
CREATE INDEX "idx_divisions_deleted_at" 
ON "divisions" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 5. Trigger updated_at Otomatis
CREATE TRIGGER update_divisions_modtime
    BEFORE UPDATE ON "divisions"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();