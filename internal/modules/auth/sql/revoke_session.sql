UPDATE user_sessions
SET status = 'revoked', revoked_at = now(), updated_at = now()
WHERE refresh_token = $1 AND deleted_at IS NULL
