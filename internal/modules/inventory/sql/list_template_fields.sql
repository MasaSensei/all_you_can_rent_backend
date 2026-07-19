SELECT id, tenant_id, asset_template_id, field_name, field_label, field_type,
       is_required, default_value, options, sort_order, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM template_fields
WHERE asset_template_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY sort_order ASC
