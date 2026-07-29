SELECT subtotal, discount_total
FROM bookings
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
