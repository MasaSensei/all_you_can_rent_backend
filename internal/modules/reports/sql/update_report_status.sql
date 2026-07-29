UPDATE reports
SET status     = $3,
    file_url   = COALESCE($4, file_url),
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
