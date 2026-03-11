-- 1. Hapus Indeks
DROP INDEX IF EXISTS "idx_rt_expired_cleanup";
DROP INDEX IF EXISTS "idx_rt_device_lookup";
DROP INDEX IF EXISTS "idx_rt_user_active_sessions";
DROP INDEX IF EXISTS "idx_rt_valid_token_lookup";

-- 2. Hapus Foreign Keys
ALTER TABLE "refresh_tokens" DROP CONSTRAINT IF EXISTS "fk_rt_replaced_by";
ALTER TABLE "refresh_tokens" DROP CONSTRAINT IF EXISTS "fk_rt_device";
ALTER TABLE "refresh_tokens" DROP CONSTRAINT IF EXISTS "fk_rt_user";

-- 3. Hapus Tabel
DROP TABLE IF EXISTS "refresh_tokens";