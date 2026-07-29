UPDATE notifications
SET is_read = true, read_at = now(), updated_at = now(), version = version + 1
WHERE tenant_id = $1
  AND (user_id = $2 OR customer_id = $3)
  AND is_read = false
  AND deleted_at IS NULL
