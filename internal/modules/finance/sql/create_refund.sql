INSERT INTO refunds (
    id, tenant_id, payment_id, amount, reason,
    refund_status, processed_at, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    'processed', now(), 'active',
    $6, $6, now(), now(), 1
)
