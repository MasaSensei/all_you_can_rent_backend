UPDATE menu_items SET deleted_at = now(), deleted_by = $3, status = 'deleted',
    updated_at = now(), version = version + 1
WHERE id = $1 AND menu_id = $2 AND deleted_at IS NULL
