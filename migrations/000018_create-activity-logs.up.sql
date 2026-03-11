CREATE TABLE IF NOT EXISTS "activity_logs" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid NOT NULL, -- Diwajibkan agar filtering tenant selalu aman
    "user_id" uuid,
    "action" text NOT NULL, -- Contoh: 'PROJECT_CREATED', 'INVOICE_PAID'
    "entity_type" text NOT NULL, -- Contoh: 'PROJECT', 'INVOICE'
    "entity_id" uuid,
    "meta" jsonb DEFAULT '{}'::jsonb NOT NULL, -- Menyimpan konteks seperti nama entitas agar tidak perlu JOIN berat
    "created_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys
-- Activity feed biasanya ikut dihapus jika perusahaan dihapus (Cascade)
-- Namun tetap ada jika user hanya dihapus (Set Null) agar riwayat perusahaan tetap utuh
ALTER TABLE "activity_logs" 
ADD CONSTRAINT "fk_activity_logs_company" FOREIGN KEY ("company_id") REFERENCES "companies"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_activity_logs_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL;

-- 2. Optimized Index untuk Company Activity Feed
-- Query: "Tampilkan 20 aktivitas terbaru di perusahaan ini"
CREATE INDEX "idx_activity_logs_company_feed" 
ON "activity_logs" ("company_id", "created_at" DESC);

-- 3. Optimized Index untuk User Activity Feed
-- Query: "Tampilkan aktivitas saya baru-baru ini"
CREATE INDEX "idx_activity_logs_user_feed" 
ON "activity_logs" ("user_id", "created_at" DESC) 
WHERE "user_id" IS NOT NULL;

-- 4. Optimized Index untuk Entity History
-- Query: "Tampilkan riwayat komentar/aktivitas pada Invoice ID X"
CREATE INDEX "idx_activity_logs_entity_history" 
ON "activity_logs" ("entity_type", "entity_id", "created_at" DESC);

-- 5. Gin Index untuk Pencarian di Meta (Opsional)
-- Berguna jika Anda ingin mencari kata kunci di dalam JSONB meta
CREATE INDEX "idx_activity_logs_meta_gin" ON "activity_logs" USING gin ("meta");