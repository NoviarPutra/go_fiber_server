BEGIN;

-- Hapus data dengan urutan terbalik dari Foreign Key
DELETE FROM "user_roles" WHERE "company_user_id" = '88888888-8888-8888-8888-888888888888';
DELETE FROM "company_users" WHERE "id" = '88888888-8888-8888-8888-888888888888';
DELETE FROM "user_profiles" WHERE "user_id" = '99999999-9999-9999-9999-999999999999';
DELETE FROM "users" WHERE "id" = '99999999-9999-9999-9999-999999999999';
DELETE FROM "role_permissions" WHERE "role_id" IN ('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222');
DELETE FROM "roles" WHERE "company_id" = '00000000-0000-0000-0000-000000000000';
DELETE FROM "permissions" WHERE "code" IN ('system.manage', 'billing.manage', 'user.create', 'user.read', 'user.update', 'user.delete', 'role.manage', 'audit.view', 'division.manage', 'team.manage');
DELETE FROM "companies" WHERE "id" = '00000000-0000-0000-0000-000000000000';

COMMIT;