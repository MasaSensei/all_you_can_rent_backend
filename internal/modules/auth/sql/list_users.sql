SELECT id, tenant_id, email, username, password_hash, first_name, last_name,
       phone, avatar_url, is_active, last_login_at, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM users
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
