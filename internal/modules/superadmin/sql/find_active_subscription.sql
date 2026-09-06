SELECT * FROM tenant_subscriptions
WHERE tenant_id = $1 AND status = 'active' AND deleted_at IS NULL
ORDER BY created_at DESC LIMIT 1
