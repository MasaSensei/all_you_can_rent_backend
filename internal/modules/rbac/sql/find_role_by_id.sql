SELECT id, tenant_id, name, description, is_system, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM roles
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
