UPDATE webhooks
SET url       = COALESCE($3, url),
    events    = COALESCE($4, events),
    is_active = COALESCE($5, is_active),
    updated_by = $2,
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
