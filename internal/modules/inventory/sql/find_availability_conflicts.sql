SELECT COUNT(*)
FROM asset_availability
WHERE asset_id = $1
  AND deleted_at IS NULL
  AND status = 'active'
  AND start_date < $3
  AND end_date > $2
