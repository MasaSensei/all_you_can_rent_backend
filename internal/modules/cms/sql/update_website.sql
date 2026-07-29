UPDATE websites SET title = COALESCE($3, title), theme = COALESCE($4, theme),
    updated_by = $2, updated_at = now(), version = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
