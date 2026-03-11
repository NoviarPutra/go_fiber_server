CREATE TABLE IF NOT EXISTS "user_profiles" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "user_id" uuid NOT NULL,
    "full_name" text,
    "phone" text,
    "avatar_url" text,
    "bio" text,
    "address" text,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    
    -- Constraint tambahan untuk keamanan data
    CONSTRAINT "user_profiles_user_id_unique" UNIQUE("user_id"),
    CONSTRAINT "user_profiles_phone_check" CHECK (char_length(phone) >= 7 OR phone IS NULL)
);

-- 1. Foreign Key dengan Cascade Delete
-- Saat user dihapus permanen, profil otomatis hilang (mencegah data yatim/orphan)
ALTER TABLE "user_profiles" 
ADD CONSTRAINT "fk_user_profiles_user" 
FOREIGN KEY ("user_id") REFERENCES "users"("id") 
ON DELETE CASCADE ON UPDATE NO ACTION;

-- 2. Index untuk Look-up Cepat
-- Meskipun sudah ada UNIQUE constraint (yang otomatis membuat index), 
-- kita pastikan performa JOIN users -> user_profiles tetap optimal.
CREATE INDEX "idx_user_profiles_user_id" ON "user_profiles" USING btree ("user_id");

-- 3. Partial Index untuk Phone (Opsional tapi membantu)
-- Jika Anda sering mencari user berdasarkan nomor telepon
CREATE INDEX "idx_user_profiles_phone_search" 
ON "user_profiles" ("phone") 
WHERE "phone" IS NOT NULL;

-- 4. Trigger updated_at Otomatis
-- Karena user_profiles sering diupdate (ganti bio/avatar), kita gunakan trigger yang sama
CREATE TRIGGER update_user_profiles_modtime
    BEFORE UPDATE ON "user_profiles"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();