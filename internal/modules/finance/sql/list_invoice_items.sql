SELECT id, tenant_id, invoice_id, booking_item_id,
       description, quantity, unit_price, tax_amount, line_total,
       status, created_at, updated_at, deleted_at, version
FROM invoice_items
WHERE invoice_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY created_at ASC
