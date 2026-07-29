SELECT id, tenant_id, customer_id, address_type, line1, line2,
       city, state, postal_code, country, is_default, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM customer_addresses
WHERE customer_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY is_default DESC, created_at ASC
