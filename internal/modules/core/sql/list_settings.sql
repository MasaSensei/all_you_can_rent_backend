SELECT id, tenant_id, key, value, type, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM settings
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY key ASC
