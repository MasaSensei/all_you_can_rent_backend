SELECT id, tenant_id, parent_id, name, slug, description, icon, sort_order, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM categories
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC, name ASC
