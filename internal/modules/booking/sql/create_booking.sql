INSERT INTO bookings (
    id, tenant_id, customer_id, coupon_id, booking_number,
    start_date, end_date, subtotal, tax_total, discount_total, total_amount,
    booking_status, payment_status, notes, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11,
    'pending', 'unpaid', $12, 'active',
    $13, $13, now(), now(), 1
)
