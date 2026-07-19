INSERT INTO booking_items (
    id, tenant_id, booking_id, asset_id, quantity,
    unit_price, line_total, start_date, end_date, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, 'active',
    $10, $10, now(), now(), 1
)
