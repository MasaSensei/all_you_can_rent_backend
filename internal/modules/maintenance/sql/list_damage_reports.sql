SELECT id, tenant_id, asset_id, booking_id, inspection_id,
       description, severity, repair_cost, charged_amount,
       report_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM damage_reports
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR asset_id = $2)
  AND ($3::varchar IS NULL OR report_status = $3)
  AND ($4::varchar IS NULL OR severity = $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6
