SELECT DISTINCT p.name
FROM permissions p
INNER JOIN role_permissions rp ON rp.permission_id = p.id
INNER JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = $1 AND ur.tenant_id = $2
  AND ur.deleted_at IS NULL
  AND rp.deleted_at IS NULL
  AND p.deleted_at IS NULL
