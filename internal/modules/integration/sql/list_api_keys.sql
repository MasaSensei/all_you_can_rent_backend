SELECT id, tenant_id, name, key_prefix, key_hash, scopes,
       last_used_at, expires_at, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM api_keys
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
