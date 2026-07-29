UPDATE api_keys
SET status = 'revoked', deleted_at = now(), deleted_by = $2,
    updated_at = now(), version = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
