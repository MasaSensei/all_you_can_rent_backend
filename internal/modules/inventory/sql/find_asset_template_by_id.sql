SELECT id, tenant_id, category_id, name, description, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM asset_templates
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
