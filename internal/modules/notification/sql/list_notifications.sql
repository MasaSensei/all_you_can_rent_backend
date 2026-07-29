SELECT id, tenant_id, user_id, customer_id, channel,
       title, message, is_read, read_at, status,
       created_at, updated_at, deleted_at, version
FROM notifications
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND (user_id = $2 OR customer_id = $3)
  AND ($4::boolean IS NULL OR is_read = $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6
