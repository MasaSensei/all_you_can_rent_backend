SELECT id, tenant_id, name, description, points_per_currency, redemption_rate,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM loyalty_programs
WHERE tenant_id = $1 AND deleted_at IS NULL AND status = 'active'
ORDER BY name ASC
