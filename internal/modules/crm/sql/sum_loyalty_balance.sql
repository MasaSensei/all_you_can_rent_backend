SELECT COALESCE(SUM(
    CASE
        WHEN transaction_type IN ('earn', 'adjust') THEN points
        WHEN transaction_type IN ('redeem', 'expire') THEN -points
        ELSE 0
    END
), 0)
FROM loyalty_transactions
WHERE customer_id = $1
  AND loyalty_program_id = $2
  AND deleted_at IS NULL
