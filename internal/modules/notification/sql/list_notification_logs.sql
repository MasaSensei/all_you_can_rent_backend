SELECT id, tenant_id, notification_template_id, recipient_id, recipient_type,
       channel, delivery_status, error_message, sent_at, status,
       created_at, updated_at, deleted_at, version
FROM notification_logs
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
