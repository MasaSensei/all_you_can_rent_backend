SELECT id, tenant_id, name, channel, subject, body, event_trigger,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM notification_templates
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
