SELECT id, tenant_id, asset_id, start_date, end_date, availability_type, reason, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM asset_availability
WHERE asset_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY start_date ASC
