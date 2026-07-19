INSERT INTO asset_templates (id, tenant_id, category_id, name, description, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, 'active', $6, $6, now(), now(), 1)
