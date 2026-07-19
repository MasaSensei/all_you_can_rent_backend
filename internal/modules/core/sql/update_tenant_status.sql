UPDATE tenants
SET status = $2, updated_at = now(), version = version + 1
WHERE id = $1 AND deleted_at IS NULL
