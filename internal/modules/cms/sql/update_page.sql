UPDATE pages SET title = COALESCE($3, title), content = COALESCE($4, content),
    template = COALESCE($5, template), status = COALESCE($6, status),
    updated_by = $2, updated_at = now(), version = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
