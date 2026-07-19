SELECT p.id, p.name, p.module, p.action, p.description, p.status,
       p.created_at, p.updated_at, p.deleted_at, p.version
FROM permissions p
INNER JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1 AND rp.tenant_id = $2
  AND rp.deleted_at IS NULL AND p.deleted_at IS NULL
ORDER BY p.module, p.action
