INSERT INTO damage_reports (
    id, tenant_id, asset_id, booking_id, inspection_id,
    description, severity, repair_cost, charged_amount,
    report_status, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    'open', 'active',
    $10, $10, now(), now(), 1
)
