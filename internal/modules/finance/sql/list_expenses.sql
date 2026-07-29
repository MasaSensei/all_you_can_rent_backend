SELECT id, tenant_id, asset_id, category, amount, expense_date,
       description, vendor, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM expenses
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::uuid IS NULL OR asset_id = $2)
  AND ($3::varchar IS NULL OR category = $3)
ORDER BY expense_date DESC
LIMIT $4 OFFSET $5
