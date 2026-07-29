INSERT INTO promotions (
    id, tenant_id, name, description, promotion_type, value,
    start_date, end_date, status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, 'active', $9, $9, now(), now(), 1
)
