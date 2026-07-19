UPDATE role_permissions
SET deleted_at = now(), status = 'deleted', updated_at = now()
WHERE role_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
