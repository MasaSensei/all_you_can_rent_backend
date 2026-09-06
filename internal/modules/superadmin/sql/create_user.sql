INSERT INTO users (id, tenant_id, username, email, password_hash,
                   first_name, last_name, is_active, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
