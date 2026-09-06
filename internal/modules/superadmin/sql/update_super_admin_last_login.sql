UPDATE super_admins SET last_login_at = now(), updated_at = now()
WHERE id = $1
