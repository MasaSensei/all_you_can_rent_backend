INSERT INTO pricing_rules (
    id, tenant_id, category_id, asset_id, name, rule_type, value,
    duration_unit, min_duration, max_duration, valid_from, valid_to,
    status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12,
    'active', $13, $13, now(), now(), 1
)
