INSERT INTO memberships (
    id, tenant_id, customer_id, plan_name, tier,
    start_date, end_date, fee, membership_status, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, 'active', 'active',
    $9, $9, now(), now(), 1
)
