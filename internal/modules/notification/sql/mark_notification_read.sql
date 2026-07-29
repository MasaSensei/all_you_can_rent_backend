UPDATE notifications
SET is_read = true, read_at = now(), updated_at = now(), version = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
