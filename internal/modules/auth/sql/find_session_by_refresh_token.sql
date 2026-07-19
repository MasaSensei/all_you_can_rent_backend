SELECT id, tenant_id, user_id, token, refresh_token,
       ip_address, user_agent, expires_at, status,
       created_at, updated_at, deleted_at, version
FROM user_sessions
WHERE refresh_token = $1 AND status = 'active' AND deleted_at IS NULL
  AND expires_at > now()
