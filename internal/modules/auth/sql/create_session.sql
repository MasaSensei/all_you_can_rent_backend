INSERT INTO user_sessions (
    id, tenant_id, user_id, token, refresh_token,
    ip_address, user_agent, expires_at, status, created_at, updated_at, version
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active', now(), now(), 1)
