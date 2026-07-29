SELECT id, tenant_id, name, key_prefix, key_hash, scopes,
       last_used_at, expires_at, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM api_keys
WHERE key_hash = $1 AND status = 'active' AND deleted_at IS NULL
  AND (expires_at IS NULL OR expires_at > now())
