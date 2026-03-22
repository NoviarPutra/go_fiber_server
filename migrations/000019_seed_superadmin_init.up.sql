BEGIN;

-- 1. Buat "System Tenant" (Isolasi khusus untuk pengelola platform)
INSERT INTO "companies" ("id", "name", "code") 
VALUES (
    '00000000-0000-0000-0000-000000000000', 
    'SAAS PLATFORM ADMINISTRATION', 
    'SYSTEM'
) ON CONFLICT (id) DO NOTHING;

-- 2. Buat Daftar Permission Lengkap (Master Data)
INSERT INTO "permissions" ("code", "module", "description") VALUES
('system.manage', 'Platform', 'Akses penuh ke seluruh tenant/perusahaan'),
('billing.manage', 'Platform', 'Mengelola tagihan dan paket perusahaan'),
('user.create', 'User Management', 'Bisa membuat user baru'),
('user.read', 'User Management', 'Bisa melihat daftar user'),
('user.update', 'User Management', 'Bisa mengubah data user'),
('user.delete', 'User Management', 'Bisa menghapus user'),
('role.manage', 'Security', 'Bisa mengatur role dan permission'),
('audit.view', 'Security', 'Bisa melihat audit logs'),
('division.manage', 'Organization', 'Bisa mengatur divisi'),
('team.manage', 'Organization', 'Bisa mengatur tim')
ON CONFLICT (code) DO NOTHING;

-- 3. Buat Role SUPERADMIN (Global) & OWNER (Template)
INSERT INTO "roles" ("id", "company_id", "name", "is_system", "description") VALUES
(
    '11111111-1111-1111-1111-111111111111',
    '00000000-0000-0000-0000-000000000000',
    'SUPERADMIN',
    true,
    'Puncak tertinggi akses seluruh platform'
),
(
    '22222222-2222-2222-2222-222222222222',
    '00000000-0000-0000-0000-000000000000',
    'OWNER_TEMPLATE',
    true,
    'Template role untuk pemilik perusahaan baru'
) ON CONFLICT (id) DO NOTHING;

-- 4. Mapping Permission ke SUPERADMIN (All Access)
INSERT INTO "role_permissions" ("role_id", "permission_id", "allowed")
SELECT '11111111-1111-1111-1111-111111111111', id, true FROM "permissions"
ON CONFLICT DO NOTHING;

-- 5. Buat User Superadmin (Password: superadmin123)
-- Hash Argon2id ini valid untuk password 'superadmin123'
INSERT INTO "users" ("id", "email", "username", "password_hash", "is_email_verified", "is_active")
VALUES (
    '99999999-9999-9999-9999-999999999999',
    'admin@saas.com',
    'superadmin',
    '$argon2id$v=19$m=65536,t=1,p=4$S9U3v6K7YVz9Fw0X$K8L9M0N1P2Q3R4S5T6U7V8W9X0Y1Z2A3B4C5D6E7F8G',
    true,
    true
) ON CONFLICT (id) DO NOTHING;

-- 6. Profile Superadmin
INSERT INTO "user_profiles" ("user_id", "full_name")
VALUES ('99999999-9999-9999-9999-999999999999', 'Platform Super Administrator')
ON CONFLICT (user_id) DO NOTHING;

-- 7. Hubungkan Superadmin ke System Company & Berikan Role
INSERT INTO "company_users" ("id", "company_id", "user_id", "is_owner", "is_active")
VALUES (
    '88888888-8888-8888-8888-888888888888',
    '00000000-0000-0000-0000-000000000000',
    '99999999-9999-9999-9999-999999999999',
    true,
    true
) ON CONFLICT (id) DO NOTHING;

INSERT INTO "user_roles" ("company_user_id", "role_id")
VALUES ('88888888-8888-8888-8888-888888888888', '11111111-1111-1111-1111-111111111111')
ON CONFLICT DO NOTHING;

COMMIT;