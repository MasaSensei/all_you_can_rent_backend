SELECT id, tenant_id, payment_id, amount, reason,
       refund_status, processed_at, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM refunds
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
