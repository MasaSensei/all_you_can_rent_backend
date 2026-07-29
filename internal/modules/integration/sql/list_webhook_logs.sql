SELECT id, tenant_id, webhook_id, event_type, payload,
       response_code, response_body, triggered_at, status,
       created_at, updated_at, deleted_at, version
FROM webhook_logs
WHERE webhook_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY triggered_at DESC
LIMIT $3 OFFSET $4
