INSERT INTO expenses (
    id, tenant_id, asset_id, category, amount, expense_date,
    description, vendor, status,
    created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, 'active',
    $9, $9, now(), now(), 1
)
