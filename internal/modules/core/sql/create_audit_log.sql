INSERT INTO audit_logs
    (id, tenant_id, user_id, entity_type, entity_id, action, old_values, new_values,
     ip_address, user_agent, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active', now(), now(), 1)
