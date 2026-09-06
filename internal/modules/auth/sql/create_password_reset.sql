INSERT INTO password_resets (id, tenant_id, user_id, token_hash, expires_at, status)
VALUES ($1, $2, $3, $4, $5, 'active')
