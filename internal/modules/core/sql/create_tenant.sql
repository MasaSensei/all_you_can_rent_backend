INSERT INTO tenants (id, name, slug, domain, plan, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, now(), now(), 1)
