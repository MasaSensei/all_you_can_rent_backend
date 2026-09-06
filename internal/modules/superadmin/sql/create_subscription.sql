INSERT INTO tenant_subscriptions
  (id, tenant_id, plan_id, billing_cycle, price_paid, starts_at, ends_at, auto_renew, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
