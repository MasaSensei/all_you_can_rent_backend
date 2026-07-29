SELECT COALESCE(SUM(amount), 0)
FROM refunds
WHERE payment_id = $1
  AND refund_status = 'processed'
  AND deleted_at IS NULL
