INSERT INTO payments (
    id, tenant_id, invoice_id, customer_id, payment_method,
    transaction_reference, amount, currency, paid_at, payment_status,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, now(), 'succeeded',
    'active', $9, $9, now(), now(), 1
)
