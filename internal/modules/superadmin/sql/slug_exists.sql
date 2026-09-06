SELECT EXISTS(SELECT 1 FROM tenants WHERE slug = $1 AND deleted_at IS NULL)
