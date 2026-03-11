CREATE TABLE IF NOT EXISTS "audit_logs" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_id" uuid,
    "user_id" uuid,
    "action" text NOT NULL, -- Contoh: 'INSERT', 'UPDATE', 'DELETE', 'LOGIN'
    "table_name" text NOT NULL,
    "record_id" uuid, -- ID dari data yang diubah
    "old_data" jsonb,
    "new_data" jsonb,
    "ip_address" text,
    "user_agent" text,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL
);

-- 1. Foreign Keys dengan "SET NULL"
-- Log audit harus tetap ada sebagai bukti sejarah meskipun perusahaan atau user dihapus
ALTER TABLE "audit_logs" 
ADD CONSTRAINT "fk_audit_logs_company" FOREIGN KEY ("company_id") REFERENCES "companies"("id") ON DELETE SET NULL,
ADD CONSTRAINT "fk_audit_logs_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL;

-- 2. Optimized Index untuk Admin Panel
-- Query: "Tampilkan semua perubahan yang terjadi di perusahaan X"
CREATE INDEX "idx_audit_logs_company_lookup" 
ON "audit_logs" ("company_id", "created_at" DESC) 
WHERE "company_id" IS NOT NULL;

-- 3. Optimized Index untuk User Activity
-- Query: "Tampilkan riwayat aktifitas user Y"
CREATE INDEX "idx_audit_logs_user_lookup" 
ON "audit_logs" ("user_id", "created_at" DESC) 
WHERE "user_id" IS NOT NULL;

-- 4. Optimized Index untuk Record History
-- Query: "Siapa saja yang pernah mengubah data produk dengan ID Z?"
CREATE INDEX "idx_audit_logs_record_lookup" 
ON "audit_logs" ("table_name", "record_id", "created_at" DESC);

-- 5. Time-based Index untuk Maintenance
-- Berguna untuk archiving atau pembersihan log lama (misal log > 1 tahun)
CREATE INDEX "idx_audit_logs_created_at_desc" 
ON "audit_logs" ("created_at" DESC);