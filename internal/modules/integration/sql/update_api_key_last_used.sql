UPDATE api_keys
SET last_used_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
