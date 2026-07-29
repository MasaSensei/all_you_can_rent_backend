SELECT id, tenant_id, name, report_type, parameters, generated_format,
       file_url, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM reports
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
