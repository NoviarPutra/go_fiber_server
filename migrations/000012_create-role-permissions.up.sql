CREATE TABLE IF NOT EXISTS "role_permissions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "role_id" uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    "allowed" boolean DEFAULT true NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys dengan Integritas Tinggi
-- Jika Role atau Permission dihapus, relasi di tabel ini otomatis hilang (Cascade)
ALTER TABLE "role_permissions" 
ADD CONSTRAINT "fk_rp_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_rp_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE;

-- 2. Strict Unique Constraint
-- Menjamin satu Role tidak bisa memiliki entri ganda untuk satu Permission yang sama
CREATE UNIQUE INDEX "idx_role_permission_unique" ON "role_permissions" ("role_id", "permission_id");

-- 3. Optimized Index untuk Reverse Lookup
-- Sangat berguna saat Anda ingin tahu: "Role mana saja yang punya akses ke permission X?"
CREATE INDEX "idx_role_permissions_permission_id" ON "role_permissions" ("permission_id");

-- 4. Optimized Index untuk Checking
-- Mempercepat filter untuk izin yang secara eksplisit diperbolehkan
CREATE INDEX "idx_role_permissions_allowed" ON "role_permissions" ("role_id") WHERE "allowed" IS TRUE;