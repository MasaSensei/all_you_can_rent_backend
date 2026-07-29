SELECT id, tenant_id, name, description, promotion_type, value,
       start_date, end_date, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM promotions
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND status = 'active'
  AND (start_date IS NULL OR start_date <= now())
  AND (end_date IS NULL OR end_date >= now())
ORDER BY created_at DESC
