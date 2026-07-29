SELECT COALESCE(SUM(p.amount), 0)::float8
FROM payments p
INNER JOIN invoices i ON i.id = p.invoice_id
WHERE i.tenant_id = $1
  AND p.payment_status = 'succeeded'
  AND p.deleted_at IS NULL
  AND p.created_at BETWEEN $2 AND $3
