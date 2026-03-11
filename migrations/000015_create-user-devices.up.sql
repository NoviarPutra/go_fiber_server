CREATE TABLE IF NOT EXISTS "user_devices" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "user_id" uuid NOT NULL,
    "device_name" text,
    "device_type" text, -- Contoh: 'ios', 'android', 'web'
    "os" text,
    "last_active" timestamp with time zone DEFAULT now() NOT NULL,
    "push_token" text,
    "is_online" boolean DEFAULT false NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "revoked_at" timestamp with time zone -- NULL berarti sesi masih aktif/valid
);

-- 1. Foreign Key dengan Cascade Delete
-- Jika User dihapus permanen, seluruh sesi perangkat otomatis bersih
ALTER TABLE "user_devices" 
ADD CONSTRAINT "fk_user_devices_user" 
FOREIGN KEY ("user_id") REFERENCES "users"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Unique Index untuk Push Token
-- Mencegah satu token terdaftar dua kali, tapi mengizinkan NULL jika user belum grant permission push
CREATE UNIQUE INDEX "idx_user_devices_push_token_unique" 
ON "user_devices" ("push_token") 
WHERE "push_token" IS NOT NULL AND "revoked_at" IS NULL;

-- 3. Optimized Index untuk Sesi Aktif User
-- Sangat kencang untuk query: "Tampilkan perangkat saya yang sedang login"
CREATE INDEX "idx_user_devices_user_active_sessions" 
ON "user_devices" ("user_id") 
WHERE "revoked_at" IS NULL;

-- 4. Index untuk Real-time Monitoring
-- Digunakan untuk fitur 'Who is online' atau sinkronisasi WebSocket
CREATE INDEX "idx_user_devices_online_status" 
ON "user_devices" ("is_online") 
WHERE "is_online" IS TRUE AND "revoked_at" IS NULL;

-- 5. Trigger updated_at Otomatis
CREATE TRIGGER update_user_devices_modtime
    BEFORE UPDATE ON "user_devices"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();