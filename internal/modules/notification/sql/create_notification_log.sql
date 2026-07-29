INSERT INTO notification_logs (
    id, tenant_id, notification_template_id, recipient_id, recipient_type,
    channel, delivery_status, error_message, sent_at, status,
    created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, 'active',
    now(), now(), 1
)
