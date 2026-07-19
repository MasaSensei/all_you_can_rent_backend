SELECT id, tenant_id, customer_id, coupon_id, booking_number,
       start_date, end_date, subtotal, tax_total, discount_total, total_amount,
       booking_status, payment_status, notes, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM bookings
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
