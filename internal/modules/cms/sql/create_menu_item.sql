INSERT INTO menu_items (id, tenant_id, menu_id, parent_id, label, url, sort_order, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', $8, $8, now(), now(), 1)
