SELECT id, tenant_id, loyalty_program_id, customer_id, booking_id,
       points, transaction_type, description, status,
       created_at, updated_at, deleted_at, version
FROM loyalty_transactions
WHERE customer_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
