INSERT INTO invoice_items (
    id, tenant_id, invoice_id, booking_item_id,
    description, quantity, unit_price, tax_amount, line_total,
    status, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9,
    'active', now(), now(), 1
)
