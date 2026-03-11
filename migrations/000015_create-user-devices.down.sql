-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_user_devices_modtime ON "user_devices";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_user_devices_online_status";
DROP INDEX IF EXISTS "idx_user_devices_user_active_sessions";
DROP INDEX IF EXISTS "idx_user_devices_push_token_unique";

-- 3. Hapus Foreign Key
ALTER TABLE "user_devices" DROP CONSTRAINT IF EXISTS "fk_user_devices_user";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "user_devices";