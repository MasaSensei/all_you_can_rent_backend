SELECT id, tenant_id, category_id, asset_id, name, rule_type, value,
       duration_unit, min_duration, max_duration, valid_from, valid_to,
       status, created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM pricing_rules
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND status = 'active'
  AND (valid_from IS NULL OR valid_from <= now())
  AND (valid_to IS NULL OR valid_to >= now())
  AND (asset_id = $2 OR category_id = $3 OR (asset_id IS NULL AND category_id IS NULL))
ORDER BY
  CASE WHEN asset_id = $2 THEN 0
       WHEN category_id = $3 THEN 1
       ELSE 2 END ASC
LIMIT 1
