CREATE TABLE IF NOT EXISTS "permissions" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "code" text NOT NULL,
    "module" text NOT NULL, -- Diwajibkan agar pengelompokan di UI mudah
    "description" text,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,

    -- Constraint: Pastikan kode mengikuti format snake_case atau dot.notation (misal: user.create)
    CONSTRAINT "permissions_code_check" CHECK (char_length(code) >= 3),
    CONSTRAINT "permissions_code_unique" UNIQUE("code")
);

-- 1. Index untuk Lookup Cepat berdasarkan Module
-- Sangat berguna saat menampilkan daftar permission per kategori di halaman Role Management
CREATE INDEX "idx_permissions_module" ON "permissions" ("module");

-- 2. Index untuk Pencarian Berdasarkan Code
-- Meskipun UNIQUE sudah membuat index, secara eksplisit kita pastikan pencarian teks kencang
CREATE INDEX "idx_permissions_code_lookup" ON "permissions" USING btree ("code");

-- 3. Trigger updated_at Otomatis
CREATE TRIGGER update_permissions_modtime
    BEFORE UPDATE ON "permissions"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();