SELECT id, tenant_id, asset_id, booking_item_id, inspection_type,
       inspected_at, inspector_name, findings, result, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM inspections
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
