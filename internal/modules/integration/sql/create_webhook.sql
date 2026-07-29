INSERT INTO webhooks (
    id, tenant_id, url, events, secret, is_active,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, true,
    'active', $6, $6, now(), now(), 1
)
