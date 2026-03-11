CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "email" text NOT NULL,
    "username" text NOT NULL,
    "password_hash" text NOT NULL,
    "is_email_verified" boolean DEFAULT false NOT NULL,
    "is_customer" boolean DEFAULT false NOT NULL,
    "is_active" boolean DEFAULT true NOT NULL,
    "last_login_at" timestamp with time zone,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,
    
    -- Constraint format email sederhana di level DB
    CONSTRAINT "users_email_check" CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    -- Pastikan username tidak kosong atau terlalu pendek
    CONSTRAINT "users_username_check" CHECK (char_length(username) >= 3)
);

-- 1. Robust Unique Index (Hanya unik untuk user yang belum dihapus)
-- Ini memungkinkan re-registrasi dengan email/username yang sama jika akun lama sudah di-delete
CREATE UNIQUE INDEX "idx_users_email_active_unique" 
ON "users" ("email") 
WHERE "deleted_at" IS NULL;

CREATE UNIQUE INDEX "idx_users_username_active_unique" 
ON "users" ("username") 
WHERE "deleted_at" IS NULL;

-- 2. Optimized Lookup Index untuk Login
-- Menggunakan btree (default) tetapi difilter hanya untuk data aktif
CREATE INDEX "idx_users_auth_lookup" 
ON "users" ("email", "username") 
WHERE "deleted_at" IS NULL AND "is_active" = true;

-- 3. Optimized Index untuk Soft Delete
CREATE INDEX "idx_users_deleted_at" 
ON "users" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;

-- 4. Trigger updated_at Otomatis
CREATE TRIGGER update_users_modtime
    BEFORE UPDATE ON "users"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();