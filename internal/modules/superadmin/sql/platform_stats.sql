SELECT
  COUNT(*)                                                   AS total_tenants,
  COUNT(*) FILTER (WHERE status = 'active')                  AS active_tenants,
  COUNT(*) FILTER (WHERE subscription_status = 'trial')      AS trial_tenants,
  COALESCE(
    (SELECT SUM(ts.price_paid)
     FROM tenant_subscriptions ts
     WHERE ts.status = 'active'
       AND ts.billing_cycle = 'monthly'
       AND ts.deleted_at IS NULL), 0
  )                                                          AS total_mrr,
  COUNT(*) FILTER (
    WHERE created_at >= date_trunc('month', now())
  )                                                          AS new_tenants_this_month
FROM tenants
WHERE deleted_at IS NULL
