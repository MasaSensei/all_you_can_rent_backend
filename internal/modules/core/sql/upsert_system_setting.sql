INSERT INTO system_settings (id, key, value, type, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, 'active', now(), now(), 1)
ON CONFLICT (key) WHERE deleted_at IS NULL
DO UPDATE SET value = EXCLUDED.value, type = EXCLUDED.type,
              updated_at = now(), version = system_settings.version + 1
