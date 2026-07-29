INSERT INTO websites (id, tenant_id, domain, title, theme, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, 'active', $6, $6, now(), now(), 1)
