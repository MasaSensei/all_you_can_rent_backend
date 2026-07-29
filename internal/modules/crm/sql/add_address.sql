INSERT INTO customer_addresses (
    id, tenant_id, customer_id, address_type, line1, line2,
    city, state, postal_code, country, is_default, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, 'active',
    $12, $12, now(), now(), 1
)
