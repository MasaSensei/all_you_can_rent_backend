UPDATE coupons
SET used_count = used_count + 1,
    updated_at = now(),
    version    = version + 1
WHERE id = $1 AND deleted_at IS NULL
