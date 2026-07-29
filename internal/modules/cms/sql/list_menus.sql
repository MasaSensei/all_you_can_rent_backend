SELECT id, tenant_id, website_id, name, location, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM menus WHERE website_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY name ASC
