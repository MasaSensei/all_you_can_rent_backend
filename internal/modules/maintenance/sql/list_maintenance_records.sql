SELECT id, tenant_id, asset_id, maintenance_type, description, cost,
       scheduled_date, completed_date, performed_by, maintenance_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM maintenance_records
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR asset_id = $2)
  AND ($3::varchar IS NULL OR maintenance_status = $3)
ORDER BY scheduled_date ASC NULLS LAST, created_at DESC
LIMIT $4 OFFSET $5
