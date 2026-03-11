CREATE TABLE IF NOT EXISTS "user_permissions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_user_id" uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    "allowed" boolean DEFAULT true NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys dengan Integritas Tinggi
ALTER TABLE "user_permissions" 
ADD CONSTRAINT "fk_up_company_user" FOREIGN KEY ("company_user_id") REFERENCES "company_users"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_up_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE;

-- 2. Strict Unique Index
-- Menjamin satu user tidak memiliki dua baris untuk izin yang sama (mencegah ambiguitas allowed true/false)
CREATE UNIQUE INDEX "idx_user_permission_unique" ON "user_permissions" ("company_user_id", "permission_id");

-- 3. Optimized Index untuk Authorization Check
-- Sangat cepat untuk query: "Apakah user ini punya override izin X?"
CREATE INDEX "idx_user_permissions_lookup" 
ON "user_permissions" ("company_user_id", "permission_id") 
WHERE "allowed" IS TRUE;

-- 4. Trigger updated_at Otomatis
CREATE TRIGGER update_user_permissions_modtime
    BEFORE UPDATE ON "user_permissions"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();