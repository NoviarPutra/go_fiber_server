CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS "company_branches" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "name" text NOT NULL,
    "address" text,
    "timezone" text DEFAULT 'UTC' NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,
    
    -- Constraint untuk memastikan nama cabang tidak kosong
    CONSTRAINT "company_branches_name_check" CHECK (char_length(name) >= 1)
);

-- 1. Foreign Key dengan penamaan yang standar
ALTER TABLE "company_branches" 
ADD CONSTRAINT "fk_company_branches_company" 
FOREIGN KEY ("company_id") REFERENCES "companies"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Partial Unique Index (Bullet-proof)
-- Mencegah nama cabang yang sama di dalam satu perusahaan yang sama (hanya untuk data aktif)
CREATE UNIQUE INDEX "idx_company_branches_name_company_active_unique" 
ON "company_branches" ("company_id", "name") 
WHERE "deleted_at" IS NULL;

-- 3. Optimized Index untuk Foreign Key (Hanya data aktif)
-- Sangat berguna saat melakukan JOIN companies -> branches
CREATE INDEX "idx_company_branches_company_active" 
ON "company_branches" ("company_id") 
WHERE "deleted_at" IS NULL;

-- 4. Optimized Index untuk Soft Delete (Audit/Trash)
CREATE INDEX "idx_company_branches_deleted_at" 
ON "company_branches" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 5. Trigger untuk updated_at otomatis
CREATE TRIGGER update_company_branches_modtime
    BEFORE UPDATE ON "company_branches"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();