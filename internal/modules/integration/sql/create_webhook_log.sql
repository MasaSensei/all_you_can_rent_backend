INSERT INTO webhook_logs (
    id, tenant_id, webhook_id, event_type, payload,
    response_code, response_body, triggered_at, status,
    created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, now(), 'active',
    now(), now(), 1
)
