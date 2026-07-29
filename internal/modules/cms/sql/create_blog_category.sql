INSERT INTO blog_categories (id, tenant_id, name, slug, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, 'active', $5, $5, now(), now(), 1)
