INSERT INTO categories (id, tenant_id, parent_id, name, slug, description, icon, sort_order, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active', $9, $9, now(), now(), 1)
