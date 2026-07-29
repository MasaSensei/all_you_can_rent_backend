SELECT id, tenant_id, menu_id, parent_id, label, url, sort_order, status,
       created_by, updated_by, deleted_by, created_at, updated_at, deleted_at, version
FROM menu_items WHERE menu_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
ORDER BY sort_order ASC
