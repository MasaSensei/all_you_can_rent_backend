SELECT id, tenant_id, code, discount_type, discount_value,
       min_order_value, usage_limit, used_count, valid_from, valid_to,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM coupons
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
