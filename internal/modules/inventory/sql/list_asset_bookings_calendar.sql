SELECT
  b.id, b.booking_number,
  COALESCE(c.first_name || ' ' || c.last_name, c.company_name, 'Unknown') AS customer_name,
  b.start_date, b.end_date, b.status
FROM bookings b
JOIN customers c ON c.id = b.customer_id
WHERE b.tenant_id  = $1
  AND b.asset_id   = $2
  AND b.deleted_at IS NULL
  AND b.status    != 'cancelled'
  AND b.end_date  >= $3
  AND b.start_date <= $4
ORDER BY b.start_date ASC
