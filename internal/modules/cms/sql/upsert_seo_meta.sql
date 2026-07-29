INSERT INTO seo_meta (id, tenant_id, entity_type, entity_id, meta_title, meta_description, meta_keywords, og_image, canonical_url, status, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'active', now(), now(), 1)
ON CONFLICT (entity_type, entity_id) WHERE deleted_at IS NULL
DO UPDATE SET meta_title       = EXCLUDED.meta_title,
              meta_description = EXCLUDED.meta_description,
              meta_keywords    = EXCLUDED.meta_keywords,
              og_image         = EXCLUDED.og_image,
              canonical_url    = EXCLUDED.canonical_url,
              updated_at       = now(),
              version          = seo_meta.version + 1
