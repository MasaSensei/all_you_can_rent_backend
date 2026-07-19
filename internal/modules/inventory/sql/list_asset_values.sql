SELECT av.id, av.tenant_id, av.asset_id, av.template_field_id, av.value, av.status,
       av.created_at, av.updated_at, av.deleted_at, av.version
FROM asset_values av
WHERE av.asset_id = $1 AND av.tenant_id = $2 AND av.deleted_at IS NULL
