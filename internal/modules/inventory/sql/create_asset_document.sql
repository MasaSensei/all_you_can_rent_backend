INSERT INTO asset_documents (id, tenant_id, asset_id, title, file_url, file_type, file_size, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', $8, $8, now(), now(), 1)
