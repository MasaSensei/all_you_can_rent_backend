SELECT id, tenant_id, user_id, entity_type, entity_id, action, old_values, new_values,
       ip_address, user_agent, status, created_at, updated_at, deleted_at, version
FROM audit_logs
WHERE tenant_id = $1
  AND ($2::varchar IS NULL OR entity_type = $2)
  AND ($3::varchar IS NULL OR entity_id = $3)
  AND ($4::varchar IS NULL OR action = $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6
