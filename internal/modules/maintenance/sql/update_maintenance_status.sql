UPDATE maintenance_records
SET maintenance_status = $3,
    completed_date     = COALESCE($4, completed_date),
    cost               = COALESCE($5, cost),
    performed_by       = COALESCE($6, performed_by),
    updated_by         = $2,
    updated_at         = now(),
    version            = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
