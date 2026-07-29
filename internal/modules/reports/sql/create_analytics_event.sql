INSERT INTO analytics_events (
    id, tenant_id, user_id, customer_id,
    event_name, event_category, event_data, source,
    occurred_at, status, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9, 'active', now(), now(), 1
)
