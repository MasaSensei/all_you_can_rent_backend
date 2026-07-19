SELECT id, tenant_id, customer_id, coupon_id, booking_number,
       start_date, end_date, subtotal, tax_total, discount_total, total_amount,
       booking_status, payment_status, notes, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM bookings
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR customer_id = $2)
  AND ($3::varchar IS NULL OR booking_status = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5
