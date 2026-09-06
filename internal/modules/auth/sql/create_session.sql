INSERT INTO user_sessions (id, tenant_id, user_id, refresh_token, user_agent, ip_address, expires_at, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active')
