INSERT INTO user_roles (id, tenant_id, user_id, role_id, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, 'active', now(), now(), 1)
ON CONFLICT (user_id, role_id) WHERE deleted_at IS NULL DO NOTHING
