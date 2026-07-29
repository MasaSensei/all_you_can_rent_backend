SELECT id, tenant_id, name, channel, subject, body, event_trigger,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM notification_templates
WHERE tenant_id = $1 AND event_trigger = $2 AND channel = $3
  AND status = 'active' AND deleted_at IS NULL
LIMIT 1
