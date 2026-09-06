INSERT INTO asset_blocks (id, tenant_id, asset_id, start_date, end_date, reason, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *
