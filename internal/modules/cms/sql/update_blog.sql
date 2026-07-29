UPDATE blogs
SET title          = COALESCE($3, title),
    content        = COALESCE($4, content),
    featured_image = COALESCE($5, featured_image),
    status         = COALESCE($6, status),
    published_at   = CASE WHEN $6 = 'published' AND published_at IS NULL THEN now() ELSE published_at END,
    updated_by     = $2,
    updated_at     = now(),
    version        = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
