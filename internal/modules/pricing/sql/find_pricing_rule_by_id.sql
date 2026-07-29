SELECT id, tenant_id, category_id, asset_id, name, rule_type, value,
       duration_unit, min_duration, max_duration, valid_from, valid_to,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM pricing_rules
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
