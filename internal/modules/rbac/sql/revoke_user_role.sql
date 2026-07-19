UPDATE user_roles
SET deleted_at = now(), status = 'deleted', updated_at = now(), version = version + 1
WHERE user_id = $1 AND role_id = $2 AND tenant_id = $3 AND deleted_at IS NULL
