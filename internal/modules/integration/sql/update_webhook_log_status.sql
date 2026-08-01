UPDATE webhook_logs
SET delivery_status = $2,
    response_code   = $3,
    response_body   = $4,
    updated_at      = now(),
    version         = version + 1
WHERE id = $1
