INSERT INTO pages (id, tenant_id, website_id, title, slug, content, template, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'draft', $8, $8, now(), now(), 1)
