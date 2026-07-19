INSERT INTO settings (id, tenant_id, key, value, type, status, created_by, updated_by, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, 'active', $6, $6, now(), now(), 1)
ON CONFLICT (tenant_id, key) WHERE deleted_at IS NULL
DO UPDATE SET value = EXCLUDED.value, type = EXCLUDED.type, updated_by = EXCLUDED.updated_by,
              updated_at = now(), version = settings.version + 1
