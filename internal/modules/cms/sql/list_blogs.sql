SELECT id, tenant_id, website_id, author_id, blog_category_id, title, slug, content, featured_image,
       published_at, status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM blogs
WHERE tenant_id = $1 AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR website_id = $2)
  AND ($3::uuid IS NULL OR blog_category_id = $3)
  AND ($4::varchar IS NULL OR status = $4)
ORDER BY created_at DESC LIMIT $5 OFFSET $6
