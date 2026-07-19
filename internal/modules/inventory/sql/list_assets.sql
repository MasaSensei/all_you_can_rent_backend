SELECT id, tenant_id, category_id, asset_template_id, name, sku, serial_number,
       description, purchase_price, replacement_value, purchase_date, condition,
       location, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM assets
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR category_id = $2)
  AND ($3::varchar IS NULL OR status = $3)
  AND ($4::varchar IS NULL OR name ILIKE '%' || $4 || '%')
ORDER BY created_at DESC
LIMIT $5 OFFSET $6
