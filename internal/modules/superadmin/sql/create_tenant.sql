INSERT INTO tenants (id, name, slug, email, phone, timezone, locale, currency,
                     subscription_status, trial_ends_at, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
