INSERT INTO asset_availability (id, tenant_id, asset_id, start_date, end_date, availability_type, reason, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', $8, $8, now(), now(), 1)
