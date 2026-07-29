SELECT id, tenant_id, name, report_type, parameters, generated_format,
       file_url, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM reports
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
