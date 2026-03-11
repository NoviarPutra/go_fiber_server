CREATE TABLE IF NOT EXISTS "employee_profiles" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "company_user_id" uuid NOT NULL,
    "division_id" uuid,
    "team_id" uuid,
    "position_id" uuid,
    "manager_id" uuid, -- Self-reference ke employee_profiles.id
    "employment_type" text, -- Contoh: 'FULLTIME', 'CONTRACT', 'INTERN'
    "join_date" date,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,

    -- Constraint untuk mencegah manager menunjuk dirinya sendiri
    CONSTRAINT "employee_profiles_manager_check" CHECK (manager_id <> id),
    -- Constraint unik: satu company_user hanya boleh punya satu profile aktif
    CONSTRAINT "employee_profiles_company_user_unique" UNIQUE ("company_user_id")
);

-- 1. Foreign Keys dengan Integritas Tinggi
ALTER TABLE "employee_profiles" 
ADD CONSTRAINT "fk_ep_company_user" FOREIGN KEY ("company_user_id") REFERENCES "company_users"("id") ON DELETE CASCADE,
ADD CONSTRAINT "fk_ep_division" FOREIGN KEY ("division_id") REFERENCES "divisions"("id") ON DELETE SET NULL,
ADD CONSTRAINT "fk_ep_team" FOREIGN KEY ("team_id") REFERENCES "teams"("id") ON DELETE SET NULL,
ADD CONSTRAINT "fk_ep_position" FOREIGN KEY ("position_id") REFERENCES "positions"("id") ON DELETE SET NULL,
ADD CONSTRAINT "fk_ep_manager" FOREIGN KEY ("manager_id") REFERENCES "employee_profiles"("id") ON DELETE SET NULL;

-- 2. Optimized Lookup Index (Hanya data aktif)
CREATE INDEX "idx_ep_company_user_active" ON "employee_profiles" ("company_user_id") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_ep_division_active" ON "employee_profiles" ("division_id") WHERE "division_id" IS NOT NULL AND "deleted_at" IS NULL;
CREATE INDEX "idx_ep_team_active" ON "employee_profiles" ("team_id") WHERE "team_id" IS NOT NULL AND "deleted_at" IS NULL;
CREATE INDEX "idx_ep_position_active" ON "employee_profiles" ("position_id") WHERE "position_id" IS NOT NULL AND "deleted_at" IS NULL;
CREATE INDEX "idx_ep_manager_active" ON "employee_profiles" ("manager_id") WHERE "manager_id" IS NOT NULL AND "deleted_at" IS NULL;

-- 3. Index untuk Soft Delete
CREATE INDEX "idx_ep_deleted_at" ON "employee_profiles" ("deleted_at") WHERE "deleted_at" IS NOT NULL;

-- 4. Trigger updated_at
CREATE TRIGGER update_employee_profiles_modtime
    BEFORE UPDATE ON "employee_profiles"
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();