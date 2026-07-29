UPDATE notification_templates
SET name       = COALESCE($3, name),
    subject    = COALESCE($4, subject),
    body       = COALESCE($5, body),
    updated_by = $2,
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
