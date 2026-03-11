CREATE TABLE IF NOT EXISTS "user_roles" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_user_id" uuid NOT NULL,
    "role_id" uuid NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys dengan Integritas Tinggi
-- Jika akun karyawan (company_users) atau Role dihapus, relasi ini otomatis bersih
ALTER TABLE "user_roles" 
ADD CONSTRAINT "fk_ur_company_user" FOREIGN KEY ("company_user_id") REFERENCES "company_users"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_ur_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE;

-- 2. Strict Unique Constraint (Paling Penting!)
-- Mencegah seorang user memiliki role yang sama dua kali di perusahaan yang sama
CREATE UNIQUE INDEX "idx_user_role_unique" ON "user_roles" ("company_user_id", "role_id");

-- 3. Optimized Lookup Index
-- Mempercepat query: "Role apa saja yang dimiliki oleh user X?"
CREATE INDEX "idx_user_roles_company_user_id" ON "user_roles" ("company_user_id");

-- 4. Optimized Reverse Index
-- Mempercepat query: "Siapa saja user yang memiliki role 'Admin'?"
CREATE INDEX "idx_user_roles_role_id" ON "user_roles" ("role_id");