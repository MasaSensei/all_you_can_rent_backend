SELECT id, tenant_id, name, rate, tax_type, is_default, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM taxes
WHERE tenant_id = $1 AND is_default = true AND deleted_at IS NULL
LIMIT 1
