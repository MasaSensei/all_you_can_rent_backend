SELECT id, tenant_id, booking_id, asset_id, quantity,
       unit_price, line_total, start_date, end_date, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM booking_items
WHERE booking_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at ASC
