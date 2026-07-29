SELECT COUNT(DISTINCT customer_id)::int
FROM bookings
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND booking_status NOT IN ('cancelled')
  AND created_at BETWEEN $2 AND $3
