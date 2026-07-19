INSERT INTO password_resets (
    id, tenant_id, user_id, token, expires_at,
    used, status, created_at, updated_at, version
) VALUES ($1, $2, $3, $4, $5, false, 'active', now(), now(), 1)
