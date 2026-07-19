SELECT COUNT(*)
FROM booking_items bi
INNER JOIN bookings b ON b.id = bi.booking_id
WHERE bi.asset_id = $1
  AND b.booking_status NOT IN ('cancelled', 'completed')
  AND bi.deleted_at IS NULL
  AND b.deleted_at IS NULL
  AND bi.start_date < $3
  AND bi.end_date > $2
