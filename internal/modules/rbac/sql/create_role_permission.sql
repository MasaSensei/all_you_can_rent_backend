INSERT INTO role_permissions (id, tenant_id, role_id, permission_id, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, 'active', now(), now(), 1)
ON CONFLICT (role_id, permission_id) WHERE deleted_at IS NULL DO NOTHING
