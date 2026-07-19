UPDATE user_sessions
SET status = 'revoked', updated_at = now(), version = version + 1
WHERE user_id = $1 AND status = 'active' AND deleted_at IS NULL
