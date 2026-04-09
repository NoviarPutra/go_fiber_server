CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_logs (
        table_name, action, record_id, old_data, new_data, 
        user_id, ip_address, user_agent, company_id
    ) VALUES (
        TG_TABLE_NAME, 
        TG_OP, 
        COALESCE(NEW.id, OLD.id), 
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN row_to_json(OLD) ELSE NULL END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN row_to_json(NEW) ELSE NULL END,
        NULLIF(current_setting('app.audit_user_id', true), '')::uuid,
        NULLIF(current_setting('app.audit_ip_address', true), ''),
        NULLIF(current_setting('app.audit_user_agent', true), ''),
        NULLIF(current_setting('app.audit_company_id', true), '')::uuid
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to companies
CREATE TRIGGER audit_companies_trigger
AFTER INSERT OR UPDATE OR DELETE ON companies
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();

-- Apply to company_branches
CREATE TRIGGER audit_company_branches_trigger
AFTER INSERT OR UPDATE OR DELETE ON company_branches
FOR EACH ROW EXECUTE FUNCTION audit_trigger_func();
