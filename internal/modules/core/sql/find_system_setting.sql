SELECT id, key, value, type, status, created_at, updated_at, deleted_at, version
FROM system_settings
WHERE key = $1 AND deleted_at IS NULL
