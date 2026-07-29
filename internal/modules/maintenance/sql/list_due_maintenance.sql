SELECT id, tenant_id, asset_id, maintenance_type, description, cost,
       scheduled_date, completed_date, performed_by, maintenance_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM maintenance_records
WHERE maintenance_status = 'scheduled'
  AND deleted_at IS NULL
  AND scheduled_date <= $1
ORDER BY scheduled_date ASC
