INSERT INTO assets (id, tenant_id, category_id, asset_template_id, name, sku, serial_number, description, purchase_price, replacement_value, purchase_date, condition, location, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 'available', $14, $14, now(), now(), 1)
