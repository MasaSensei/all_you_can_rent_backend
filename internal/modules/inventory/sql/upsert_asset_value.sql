INSERT INTO asset_values (id, tenant_id, asset_id, template_field_id, value, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, 'active', now(), now(), 1)
ON CONFLICT (asset_id, template_field_id) WHERE deleted_at IS NULL
DO UPDATE SET value = EXCLUDED.value, updated_at = now(), version = asset_values.version + 1
