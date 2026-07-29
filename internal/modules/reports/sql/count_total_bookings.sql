SELECT COUNT(*)::int
FROM bookings
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND created_at BETWEEN $2 AND $3
