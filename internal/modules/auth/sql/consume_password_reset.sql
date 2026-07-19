UPDATE password_resets
SET used = true, updated_at = now(), version = version + 1
WHERE id = $1 AND deleted_at IS NULL
