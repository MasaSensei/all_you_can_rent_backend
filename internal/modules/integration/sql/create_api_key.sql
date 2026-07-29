INSERT INTO api_keys (
    id, tenant_id, name, key_prefix, key_hash, scopes,
    expires_at, status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, 'active', $8, $8, now(), now(), 1
)
