INSERT INTO template_fields (id, tenant_id, asset_template_id, field_name, field_label, field_type, is_required, default_value, options, sort_order, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active', $11, $11, now(), now(), 1)
