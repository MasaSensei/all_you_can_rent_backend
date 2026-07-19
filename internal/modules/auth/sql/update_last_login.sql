UPDATE users
SET last_login_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
