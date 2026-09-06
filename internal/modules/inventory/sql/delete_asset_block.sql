UPDATE asset_blocks
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
