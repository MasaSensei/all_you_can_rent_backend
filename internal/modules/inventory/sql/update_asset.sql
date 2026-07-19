UPDATE assets
SET name              = COALESCE($3, name),
    description       = COALESCE($4, description),
    condition         = COALESCE($5, condition),
    location          = COALESCE($6, location),
    purchase_price    = COALESCE($7, purchase_price),
    replacement_value = COALESCE($8, replacement_value),
    updated_by        = $2,
    updated_at        = now(),
    version           = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
