SELECT id, tenant_id, asset_id, url, alt_text, is_primary, sort_order, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM asset_images
WHERE asset_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY sort_order ASC
