UPDATE notification_logs
SET delivery_status = $2,
    error_message   = $3,
    sent_at         = CASE WHEN $2 = 'sent' THEN now() ELSE sent_at END,
    updated_at      = now(),
    version         = version + 1
WHERE id = $1
