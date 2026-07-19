UPDATE bookings
SET subtotal       = $3,
    tax_total      = $4,
    discount_total = $5,
    total_amount   = $6,
    coupon_id      = $7,
    updated_at     = now(),
    version        = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
