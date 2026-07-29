INSERT INTO blogs (id, tenant_id, website_id, author_id, blog_category_id, title, slug, content, featured_image, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'draft', $10, $10, now(), now(), 1)
