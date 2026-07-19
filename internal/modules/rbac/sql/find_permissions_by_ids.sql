SELECT id, name, module, action, description, status, created_at, updated_at, deleted_at, version
FROM permissions
WHERE id = ANY($1) AND deleted_at IS NULL
