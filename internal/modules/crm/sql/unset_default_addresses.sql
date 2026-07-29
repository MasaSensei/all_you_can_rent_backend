UPDATE customer_addresses
SET is_default = false, updated_at = now(), version = version + 1
WHERE customer_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
