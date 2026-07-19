SELECT id, tenant_id, name, description, is_system, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM roles
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY name ASC
