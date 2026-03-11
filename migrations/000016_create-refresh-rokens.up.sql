CREATE TABLE IF NOT EXISTS "refresh_tokens" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "user_id" uuid NOT NULL,
    "device_id" uuid,
    "token" text NOT NULL,
    "expires_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "revoked_at" timestamp with time zone,
    "replaced_by" uuid,
    
    CONSTRAINT "refresh_tokens_token_unique" UNIQUE ("token")
);

-- 1. Foreign Keys
ALTER TABLE "refresh_tokens" 
ADD CONSTRAINT "fk_rt_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_rt_device" FOREIGN KEY ("device_id") REFERENCES "user_devices"("id") ON DELETE SET NULL,
ADD CONSTRAINT "fk_rt_replaced_by" FOREIGN KEY ("replaced_by") REFERENCES "refresh_tokens"("id") ON DELETE SET NULL;

-- 2. Perbaikan Index Lookup
-- Kita hapus "expires_at > now()". 
-- Filter waktu akan dilakukan di level Query Go, tapi indeks tetap sangat cepat karena memfilter "revoked_at"
CREATE INDEX "idx_rt_valid_token_lookup" 
ON "refresh_tokens" ("token") 
WHERE "revoked_at" IS NULL;

-- 3. Perbaikan Index User Sessions
CREATE INDEX "idx_rt_user_active_sessions" 
ON "refresh_tokens" ("user_id") 
WHERE "revoked_at" IS NULL;

-- 4. Index Device Lookup (Tetap sama)
CREATE INDEX "idx_rt_device_lookup" 
ON "refresh_tokens" ("device_id") 
WHERE "device_id" IS NOT NULL;

-- 5. Perbaikan Cleanup Index
-- Cukup buat btree index biasa pada expires_at. 
-- Postgres sangat cepat melakukan range scan ( < ) pada kolom timestamp yang terindeks.
CREATE INDEX "idx_rt_expires_at" ON "refresh_tokens" ("expires_at");