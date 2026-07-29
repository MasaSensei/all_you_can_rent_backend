SELECT id, tenant_id, customer_id, booking_id, invoice_number,
       issue_date, due_date, subtotal, tax_total, discount_total,
       total_amount, amount_paid, amount_due, invoice_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM invoices
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR customer_id = $2)
  AND ($3::varchar IS NULL OR invoice_status = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5
