SELECT id, tenant_id, customer_id, plan_name, tier,
       start_date, end_date, fee, membership_status, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM memberships
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
