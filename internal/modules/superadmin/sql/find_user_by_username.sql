SELECT * FROM users WHERE username = $1 AND tenant_id = $2 AND deleted_at IS NULL
