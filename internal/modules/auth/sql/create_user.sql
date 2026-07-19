INSERT INTO users (
    id, tenant_id, email, password_hash, first_name, last_name,
    phone, is_active, status, created_at, updated_at, version
) VALUES ($1, $2, $3, $4, $5, $6, $7, true, 'active', now(), now(), 1)
