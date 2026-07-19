SELECT id, tenant_id, booking_id, booking_item_id,
       old_end_date, new_end_date, additional_cost, reason, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM booking_extensions
WHERE booking_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at DESC
