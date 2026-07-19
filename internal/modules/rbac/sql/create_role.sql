INSERT INTO roles (id, tenant_id, name, description, is_system, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, false, 'active', $5, $5, now(), now(), 1)
