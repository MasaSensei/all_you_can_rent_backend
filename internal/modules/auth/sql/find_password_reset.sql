SELECT pr.*, u.tenant_id FROM password_resets pr
JOIN users u ON u.id = pr.user_id
WHERE pr.token_hash = $1
  AND pr.status = 'active'
  AND pr.used_at IS NULL
  AND pr.expires_at > now()
  AND pr.deleted_at IS NULL
