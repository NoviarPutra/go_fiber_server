-- 1. Hapus Trigger
DROP TRIGGER IF EXISTS update_employee_profiles_modtime ON "employee_profiles";

-- 2. Hapus Indeks
DROP INDEX IF EXISTS "idx_ep_deleted_at";
DROP INDEX IF EXISTS "idx_ep_manager_active";
DROP INDEX IF EXISTS "idx_ep_position_active";
DROP INDEX IF EXISTS "idx_ep_team_active";
DROP INDEX IF EXISTS "idx_ep_division_active";
DROP INDEX IF EXISTS "idx_ep_company_user_active";

-- 3. Hapus Tabel (FK akan otomatis terhapus)
DROP TABLE IF EXISTS "employee_profiles";