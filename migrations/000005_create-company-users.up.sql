CREATE TABLE IF NOT EXISTS "company_users" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "branch_id" uuid,
    "employee_code" text,
    "is_owner" boolean DEFAULT false NOT NULL,
    "is_active" boolean DEFAULT true NOT NULL,
    "joined_at" timestamp with time zone DEFAULT now() NOT NULL,
    "left_at" timestamp with time zone,
    "deleted_at" timestamp with time zone,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys dengan Integritas Data
ALTER TABLE "company_users" 
ADD CONSTRAINT "fk_cu_company" FOREIGN KEY ("company_id") REFERENCES "companies"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_cu_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_cu_branch" FOREIGN KEY ("branch_id") REFERENCES "company_branches"("id") ON DELETE SET NULL;

-- 2. Robust Unique Constraint (Paling Penting!)
-- Mengizinkan user yang sama join kembali ke company yang sama jika record lama sudah di-delete
CREATE UNIQUE INDEX "idx_cu_company_user_active_unique" 
ON "company_users" ("company_id", "user_id") 
WHERE "deleted_at" IS NULL;

-- 3. Optimized Lookup Index untuk Cek Izin (RBAC)
-- Sangat cepat untuk query: "Apakah user X aktif di company Y?"
CREATE INDEX "idx_cu_auth_lookup" 
ON "company_users" ("user_id", "company_id") 
WHERE "deleted_at" IS NULL AND "is_active" = true;

-- 4. Index untuk Filter Branch
CREATE INDEX "idx_cu_branch_lookup" 
ON "company_users" ("branch_id") 
WHERE "branch_id" IS NOT NULL AND "deleted_at" IS NULL;

-- 5. Optimized Index untuk Soft Delete
CREATE INDEX "idx_cu_deleted_at" 
ON "company_users" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 6. Trigger updated_at
CREATE TRIGGER update_company_users_modtime
    BEFORE UPDATE ON "company_users"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();