INSERT INTO inspections (
    id, tenant_id, asset_id, booking_item_id, inspection_type,
    inspected_at, inspector_name, findings, result, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    now(), $6, $7, $8, 'active',
    $9, $9, now(), now(), 1
)
