INSERT INTO reports (
    id, tenant_id, name, report_type, parameters, generated_format,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    'queued', $7, $7, now(), now(), 1
)
