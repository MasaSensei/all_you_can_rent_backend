UPDATE users
SET first_name = COALESCE($3, first_name),
    last_name  = COALESCE($4, last_name),
    phone      = COALESCE($5, phone),
    updated_by = $2,
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
