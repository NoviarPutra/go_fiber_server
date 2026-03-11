-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_company_branches_modtime ON "company_branches";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_company_branches_deleted_at";
DROP INDEX IF EXISTS "idx_company_branches_company_active";
DROP INDEX IF EXISTS "idx_company_branches_name_company_active_unique";

-- 3. Hapus Foreign Key (Opsional karena tabel akan dihapus, tapi baik untuk dokumentasi)
ALTER TABLE "company_branches" DROP CONSTRAINT IF EXISTS "fk_company_branches_company";

-- 4. Hapus Tabel
DROP TABLE IF EXISTS "company_branches";