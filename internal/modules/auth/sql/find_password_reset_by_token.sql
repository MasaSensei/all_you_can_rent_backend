SELECT id, tenant_id, user_id, token, expires_at, used, status,
       created_at, updated_at, deleted_at, version
FROM password_resets
WHERE token = $1 AND used = false AND deleted_at IS NULL
  AND expires_at > now()
