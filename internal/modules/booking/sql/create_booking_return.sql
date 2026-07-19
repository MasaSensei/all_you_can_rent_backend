INSERT INTO booking_returns (
    id, tenant_id, booking_id, booking_item_id,
    returned_at, condition_on_return, late_fee, damage_fee, notes, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4,
    now(), $5, $6, $7, $8, 'active',
    $9, $9, now(), now(), 1
)
