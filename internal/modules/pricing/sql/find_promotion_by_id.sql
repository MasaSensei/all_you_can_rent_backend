SELECT id, tenant_id, name, description, promotion_type, value,
       start_date, end_date, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM promotions
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
