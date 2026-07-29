SELECT id, tenant_id, asset_id, booking_item_id, inspection_type,
       inspected_at, inspector_name, findings, result, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM inspections
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR asset_id = $2)
  AND ($3::varchar IS NULL OR inspection_type = $3)
  AND ($4::varchar IS NULL OR result = $4)
ORDER BY inspected_at DESC
LIMIT $5 OFFSET $6
