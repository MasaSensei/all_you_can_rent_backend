SELECT id, tenant_id, website_id, title, slug, content, template, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM pages WHERE website_id = $1 AND slug = $2 AND deleted_at IS NULL
