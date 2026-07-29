SELECT id, tenant_id, website_id, author_id, blog_category_id, title, slug, content, featured_image,
       published_at, status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM blogs WHERE website_id = $1 AND slug = $2 AND deleted_at IS NULL
