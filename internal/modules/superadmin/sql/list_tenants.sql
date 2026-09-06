SELECT * FROM tenants
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
  AND ($2 = '' OR subscription_status = $2)
  AND ($3 = '' OR status = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5
