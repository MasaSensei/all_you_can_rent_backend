SELECT id, name, slug, domain, plan, status, created_at, updated_at, deleted_at, version
FROM tenants
WHERE slug = $1 AND deleted_at IS NULL
