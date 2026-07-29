INSERT INTO loyalty_transactions (
    id, tenant_id, loyalty_program_id, customer_id, booking_id,
    points, transaction_type, description, status,
    created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, 'active',
    now(), now(), 1
)
