UPDATE categories
SET name        = COALESCE($3, name),
    description = COALESCE($4, description),
    icon        = COALESCE($5, icon),
    sort_order  = COALESCE($6, sort_order),
    updated_by  = $2,
    updated_at  = now(),
    version     = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
