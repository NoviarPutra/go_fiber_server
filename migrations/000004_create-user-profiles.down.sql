-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_user_profiles_modtime ON "user_profiles";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_user_profiles_phone_search";
DROP INDEX IF EXISTS "idx_user_profiles_user_id";

-- 3. Hapus Foreign Key
ALTER TABLE "user_profiles" DROP CONSTRAINT IF EXISTS "fk_user_profiles_user";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "user_profiles";