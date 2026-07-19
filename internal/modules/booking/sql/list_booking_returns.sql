SELECT id, tenant_id, booking_id, booking_item_id,
       returned_at, condition_on_return, late_fee, damage_fee, notes, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM booking_returns
WHERE booking_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY returned_at DESC
