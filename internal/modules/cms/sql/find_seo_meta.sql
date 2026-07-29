SELECT id, tenant_id, entity_type, entity_id, meta_title, meta_description,
       meta_keywords, og_image, canonical_url, status, created_at, updated_at, deleted_at, version
FROM seo_meta WHERE entity_type = $1 AND entity_id = $2 AND deleted_at IS NULL
