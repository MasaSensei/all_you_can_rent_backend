SELECT id, tenant_id, asset_id, title, file_url, file_type, file_size, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM asset_documents
WHERE asset_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at DESC
