INSERT INTO maintenance_records (
    id, tenant_id, asset_id, maintenance_type, description, cost,
    scheduled_date, performed_by, maintenance_status, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, 'scheduled', 'active',
    $9, $9, now(), now(), 1
)
