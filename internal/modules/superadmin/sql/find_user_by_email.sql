SELECT * FROM users
WHERE email = $1 AND tenant_id = $2 AND deleted_at IS NULL
