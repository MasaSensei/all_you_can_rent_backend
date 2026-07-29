INSERT INTO loyalty_programs (
    id, tenant_id, name, description, points_per_currency, redemption_rate,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    'active', $7, $7, now(), now(), 1
)
