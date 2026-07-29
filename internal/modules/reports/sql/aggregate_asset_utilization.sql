SELECT
    a.id   AS asset_id,
    a.name AS asset_name,
    COUNT(DISTINCT bi.id)::int                                          AS booked_days,
    GREATEST(1, ($3::date - $2::date)::int)                            AS total_days,
    ROUND(
        COUNT(DISTINCT bi.id)::numeric /
        GREATEST(1, ($3::date - $2::date)::numeric) * 100, 2
    )::float8                                                           AS utilization_pct
FROM assets a
LEFT JOIN booking_items bi ON bi.asset_id = a.id
    AND bi.deleted_at IS NULL
    AND bi.start_date >= $2
    AND bi.end_date   <= $3
WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
GROUP BY a.id, a.name
ORDER BY utilization_pct DESC
LIMIT $4
