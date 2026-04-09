DROP TRIGGER IF EXISTS audit_companies_trigger ON companies;
DROP TRIGGER IF EXISTS audit_company_branches_trigger ON company_branches;

DROP FUNCTION IF EXISTS audit_trigger_func();
