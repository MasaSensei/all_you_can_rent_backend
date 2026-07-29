SELECT id, tenant_id, url, events, secret, is_active, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM webhooks
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
