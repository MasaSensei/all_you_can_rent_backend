SELECT
    to_char(date_trunc($4, i.created_at), 'YYYY-MM-DD') AS period,
    COALESCE(SUM(p.amount), 0)::float8                  AS revenue,
    COUNT(DISTINCT i.id)::int                           AS count
FROM invoices i
LEFT JOIN payments p ON p.invoice_id = i.id
    AND p.payment_status = 'succeeded'
    AND p.deleted_at IS NULL
WHERE i.tenant_id = $1
  AND i.deleted_at IS NULL
  AND i.created_at BETWEEN $2 AND $3
GROUP BY 1
ORDER BY 1 ASC
