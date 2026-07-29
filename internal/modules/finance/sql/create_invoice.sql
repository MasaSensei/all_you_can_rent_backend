INSERT INTO invoices (
    id, tenant_id, customer_id, booking_id, invoice_number,
    issue_date, due_date, subtotal, tax_total, discount_total,
    total_amount, amount_paid, amount_due, invoice_status, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5,
    now(), $6, $7, $8, $9,
    $10, 0, $10, 'unpaid', 'active',
    $11, $11, now(), now(), 1
)
