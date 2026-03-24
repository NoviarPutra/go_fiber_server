BEGIN;

-- ============================================================
-- 1. SYSTEM TENANT (companies)
-- ============================================================
INSERT INTO "companies" ("id", "name", "code") 
VALUES (
    '50710609-8809-408a-9e6e-3f65b4c6e936', 
    'SAAS PLATFORM ADMINISTRATION', 
    'SYSTEM'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 2. BRANCH DEFAULT (company_branches)
-- ============================================================
INSERT INTO "company_branches" ("id", "company_id", "name", "address", "timezone")
VALUES (
    '00000001-0000-0000-0000-000000000001',  -- ✅ valid hex
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'Headquarters',
    'Platform HQ',
    'Asia/Jakarta'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 3. DIVISION DEFAULT (divisions)
-- ============================================================
INSERT INTO "divisions" ("id", "company_id", "name", "code")
VALUES (
    '00000002-0000-0000-0000-000000000001',  -- ✅ valid hex
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'Platform Operations',
    'PLATOPS'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 4. POSITION DEFAULT (positions)
-- ============================================================
INSERT INTO "positions" ("id", "company_id", "name", "level")
VALUES (
    '00000003-0000-0000-0000-000000000001',  -- ✅ valid hex (bukan 'p...')
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'Super Administrator',
    99
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 5. MASTER PERMISSIONS
-- ============================================================
INSERT INTO "permissions" ("code", "module", "description") VALUES
('system.manage',   'Platform',        'Akses penuh ke seluruh tenant/perusahaan'),
('billing.manage',  'Platform',        'Mengelola tagihan dan paket perusahaan'),
('user.create',     'User Management', 'Bisa membuat user baru'),
('user.read',       'User Management', 'Bisa melihat daftar user'),
('user.update',     'User Management', 'Bisa mengubah data user'),
('user.delete',     'User Management', 'Bisa menghapus user'),
('role.manage',     'Security',        'Bisa mengatur role dan permission'),
('audit.view',      'Security',        'Bisa melihat audit logs'),
('company.manage',  'Organization',    'Bisa mengatur data perusahaan'),
('branch.manage',   'Organization',    'Bisa mengatur cabang perusahaan'),
('division.manage', 'Organization',    'Bisa mengatur divisi'),
('team.manage',     'Organization',    'Bisa mengatur tim'),
('position.manage', 'Organization',    'Bisa mengatur jabatan'),
('employee.manage', 'Organization',    'Bisa mengatur data karyawan')
ON CONFLICT (code) DO NOTHING;

-- ============================================================
-- 6. ROLES: SUPERADMIN + OWNER_TEMPLATE
-- ============================================================
INSERT INTO "roles" ("id", "company_id", "name", "is_system", "description") VALUES
(
    'a7b3e942-8367-4d1a-8e2b-422998782f07',
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'SUPERADMIN',
    true,
    'Puncak tertinggi akses seluruh platform'
),
(
    'df100a94-2795-4673-8261-9f268b8b0e85',
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'OWNER_TEMPLATE',
    true,
    'Template role untuk pemilik perusahaan baru'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 7. ROLE_PERMISSIONS
-- ============================================================
INSERT INTO "role_permissions" ("role_id", "permission_id", "allowed")
SELECT 'a7b3e942-8367-4d1a-8e2b-422998782f07', id, true
FROM "permissions"
ON CONFLICT DO NOTHING;

INSERT INTO "role_permissions" ("role_id", "permission_id", "allowed")
SELECT 'df100a94-2795-4673-8261-9f268b8b0e85', id, true
FROM "permissions"
WHERE "code" NOT IN ('system.manage', 'billing.manage')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 8. USER SUPERADMIN
-- ============================================================
INSERT INTO "users" (
    "id", "email", "username", "password_hash",
    "is_email_verified", "is_active"
)
VALUES (
    'c25e4c34-8c8d-4f0e-b072-0f04c6686381',
    'admin@saas.com',
    'superadmin',
    '$argon2id$v=19$m=65536,t=1,p=8$lhtkFz6fB00gTovHXaBlvg$La9yBn9dSb4qoiyj1LcwWnvS0jtYTDlw43nrjPNB5nw',
    true,
    true
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 9. USER PROFILE
-- ============================================================
INSERT INTO "user_profiles" ("user_id", "full_name", "phone")
VALUES (
    'c25e4c34-8c8d-4f0e-b072-0f04c6686381',
    'Platform Super Administrator',
    NULL
) ON CONFLICT ("user_id") DO NOTHING;

-- ============================================================
-- 10. COMPANY_USERS
-- ============================================================
INSERT INTO "company_users" (
    "id", "company_id", "user_id", "branch_id",
    "is_owner", "is_active"
)
VALUES (
    'f92b7c0d-6e41-4b11-9a71-6789e023d8c4',
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'c25e4c34-8c8d-4f0e-b072-0f04c6686381',
    '00000001-0000-0000-0000-000000000001',  -- ← HQ branch
    true,
    true
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 11. EMPLOYEE PROFILE
-- ============================================================
INSERT INTO "employee_profiles" (
    "id", "company_user_id", "division_id", "team_id",
    "position_id", "manager_id", "employment_type", "join_date"
)
VALUES (
    '00000004-0000-0000-0000-000000000001',  -- ✅ valid hex (bukan 'e...')
    'f92b7c0d-6e41-4b11-9a71-6789e023d8c4',
    '00000002-0000-0000-0000-000000000001',
    NULL,
    '00000003-0000-0000-0000-000000000001',
    NULL,
    'FULLTIME',
    CURRENT_DATE
) ON CONFLICT ("company_user_id") DO NOTHING;

-- ============================================================
-- 12. USER_ROLES
-- ============================================================
INSERT INTO "user_roles" ("company_user_id", "role_id")
VALUES (
    'f92b7c0d-6e41-4b11-9a71-6789e023d8c4',
    'a7b3e942-8367-4d1a-8e2b-422998782f07'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 13. ACTIVITY LOG
-- ============================================================
INSERT INTO "activity_logs" (
    "company_id", "user_id",
    "action", "entity_type", "entity_id",
    "meta"
)
VALUES (
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    'c25e4c34-8c8d-4f0e-b072-0f04c6686381',
    'SYSTEM_INIT',
    'PLATFORM',
    '50710609-8809-408a-9e6e-3f65b4c6e936',
    jsonb_build_object(
        'description', 'Initial platform seed: superadmin account created',
        'version',     '1.0.0',
        'seeded_at',   now()
    )
);

COMMIT;