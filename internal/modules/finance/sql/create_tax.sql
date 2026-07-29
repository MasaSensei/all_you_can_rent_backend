INSERT INTO taxes (
    id, tenant_id, name, rate, tax_type, is_default, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6, 'active',
    $7, $7, now(), now(), 1
)
