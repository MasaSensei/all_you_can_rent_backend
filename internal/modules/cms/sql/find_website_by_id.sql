SELECT id, tenant_id, domain, title, theme, status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM websites WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
