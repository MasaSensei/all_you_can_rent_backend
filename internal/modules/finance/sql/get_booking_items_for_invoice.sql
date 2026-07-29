SELECT id, asset_id, quantity, unit_price, line_total
FROM booking_items
WHERE booking_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at ASC
