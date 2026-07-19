SELECT id, tenant_id, category_id, asset_template_id, name, sku, serial_number,
       description, purchase_price, replacement_value, purchase_date, condition,
       location, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM assets
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
