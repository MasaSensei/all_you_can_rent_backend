SELECT id, tenant_id, asset_id, booking_id, inspection_id,
       description, severity, repair_cost, charged_amount,
       report_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM damage_reports
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
