INSERT INTO coupons (
    id, tenant_id, code, discount_type, discount_value,
    min_order_value, usage_limit, used_count, valid_from, valid_to,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, 0, $8, $9,
    'active', $10, $10, now(), now(), 1
)
