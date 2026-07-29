UPDATE invoices
SET amount_paid    = $3,
    amount_due     = total_amount - $3,
    invoice_status = CASE
        WHEN $3 >= total_amount THEN 'paid'
        WHEN $3 > 0             THEN 'partial'
        ELSE 'unpaid'
    END,
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
