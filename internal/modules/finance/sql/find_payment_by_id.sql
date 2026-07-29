SELECT id, tenant_id, invoice_id, customer_id, payment_method,
       transaction_reference, amount, currency, paid_at, payment_status,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM payments
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
