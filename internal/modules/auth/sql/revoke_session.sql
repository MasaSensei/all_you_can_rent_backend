UPDATE user_sessions
SET status = 'revoked', updated_at = now(), version = version + 1
WHERE id = $1 AND deleted_at IS NULL
