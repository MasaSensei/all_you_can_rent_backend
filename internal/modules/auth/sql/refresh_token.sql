SELECT us.*, u.id AS user_id, u.tenant_id, u.email, u.password_hash,
       u.first_name, u.last_name, u.username, u.is_active
FROM user_sessions us
JOIN users u ON u.id = us.user_id
WHERE us.refresh_token = $1
  AND us.status = 'active'
  AND us.revoked_at IS NULL
  AND us.expires_at > now()
  AND us.deleted_at IS NULL
