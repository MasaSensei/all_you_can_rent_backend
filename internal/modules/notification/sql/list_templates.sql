SELECT id, tenant_id, name, channel, subject, body, event_trigger,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM notification_templates
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY event_trigger ASC, channel ASC
