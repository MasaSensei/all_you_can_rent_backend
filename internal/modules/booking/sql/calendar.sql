SELECT
  b.id,
  b.booking_number,
  COALESCE(c.first_name || ' ' || c.last_name, c.company_name, 'Unknown') AS customer_name,
  a.name       AS asset_name,
  a.id         AS asset_id,
  a.category_id,
  b.start_date,
  b.end_date,
  b.status,
  b.total_amount
FROM bookings b
JOIN assets   a ON a.id = b.asset_id
JOIN customers c ON c.id = b.customer_id
WHERE b.tenant_id = $1
  AND b.deleted_at IS NULL
  AND b.status    != 'cancelled'
  AND b.end_date  >= $2
  AND b.start_date <= $3
  AND ($4 = '' OR a.category_id::text = $4)
  AND ($5 = '' OR b.status = $5)
ORDER BY b.start_date ASC
