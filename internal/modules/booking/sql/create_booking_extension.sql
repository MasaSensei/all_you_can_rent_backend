INSERT INTO booking_extensions (
    id, tenant_id, booking_id, booking_item_id,
    old_end_date, new_end_date, additional_cost, reason, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, 'active',
    $9, $9, now(), now(), 1
)
