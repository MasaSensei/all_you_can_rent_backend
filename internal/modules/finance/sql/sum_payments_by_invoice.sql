SELECT COALESCE(SUM(amount), 0)
FROM payments
WHERE invoice_id = $1
  AND payment_status = 'succeeded'
  AND deleted_at IS NULL
