SELECT r.id, r.tenant_id, r.name, r.description, r.is_system, r.status,
       r.created_by, r.updated_by, r.deleted_by, r.created_at, r.updated_at, r.deleted_at, r.version
FROM roles r
INNER JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1 AND ur.tenant_id = $2
  AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
