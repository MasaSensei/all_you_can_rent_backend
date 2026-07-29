INSERT INTO notifications (
    id, tenant_id, user_id, customer_id, channel,
    title, message, is_read, status,
    created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, false, 'active',
    now(), now(), 1
)
