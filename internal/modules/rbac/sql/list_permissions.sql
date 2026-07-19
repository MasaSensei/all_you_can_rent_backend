SELECT id, name, module, action, description, status, created_at, updated_at, deleted_at, version
FROM permissions
WHERE deleted_at IS NULL
ORDER BY module ASC, action ASC
