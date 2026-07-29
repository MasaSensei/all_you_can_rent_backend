SELECT id, tenant_id, name, description, points_per_currency, redemption_rate,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM loyalty_programs
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
