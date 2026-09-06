SELECT * FROM asset_blocks
WHERE tenant_id = $1
  AND asset_id  = $2
  AND end_date  >= $3
  AND start_date <= $4
  AND deleted_at IS NULL
ORDER BY start_date ASC
