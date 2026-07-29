UPDATE damage_reports
SET report_status  = $3,
    repair_cost    = COALESCE($4, repair_cost),
    charged_amount = COALESCE($5, charged_amount),
    updated_by     = $2,
    updated_at     = now(),
    version        = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
