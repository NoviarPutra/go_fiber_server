CREATE TABLE IF NOT EXISTS "roles" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "name" text NOT NULL,
    "description" text,
    "is_system" boolean DEFAULT false NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,

    -- Constraint: Nama role minimal 1 karakter
    CONSTRAINT "roles_name_check" CHECK (char_length(name) >= 1)
);

-- 1. Foreign Key dengan Cascade Delete
ALTER TABLE "roles" 
ADD CONSTRAINT "fk_roles_company" 
FOREIGN KEY ("company_id") REFERENCES "companies"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Robust Unique Index untuk Nama Role
-- Mencegah nama role ganda di satu perusahaan (hanya untuk data aktif)
CREATE UNIQUE INDEX "idx_roles_name_company_active_unique" 
ON "roles" ("company_id", "name") 
WHERE "deleted_at" IS NULL;

-- 3. Optimized Index untuk Pencarian Role per Perusahaan
CREATE INDEX "idx_roles_company_active" 
ON "roles" ("company_id") 
WHERE "deleted_at" IS NULL;

-- 4. Index untuk System Roles (Guna proteksi logic di backend)
CREATE INDEX "idx_roles_system_lookup" 
ON "roles" ("is_system") 
WHERE "is_system" IS TRUE AND "deleted_at" IS NULL;

-- 5. Optimized Index untuk Soft Delete
CREATE INDEX "idx_roles_deleted_at" 
ON "roles" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 6. Trigger updated_at Otomatis
CREATE TRIGGER update_roles_modtime
    BEFORE UPDATE ON "roles"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();