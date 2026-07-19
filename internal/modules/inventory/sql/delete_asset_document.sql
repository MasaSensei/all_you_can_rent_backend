UPDATE asset_documents
SET deleted_at = now(), deleted_by = $3, status = 'deleted',
    updated_at = now(), version = version + 1
WHERE id = $1 AND asset_id = $2 AND deleted_at IS NULL
