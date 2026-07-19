SELECT id, tenant_id, email, username, password_hash, first_name, last_name,
       phone, avatar_url, is_active, last_login_at, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM users
WHERE email = $1 AND tenant_id = $2 AND deleted_at IS NULL
