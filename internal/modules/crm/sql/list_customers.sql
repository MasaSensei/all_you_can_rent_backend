SELECT id, tenant_id, first_name, last_name, email, phone, company_name,
       date_of_birth, id_document_type, id_document_number,
       customer_type, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM customers
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND ($2::varchar IS NULL OR customer_type = $2)
  AND ($3::varchar IS NULL OR (
      first_name ILIKE '%' || $3 || '%' OR
      last_name  ILIKE '%' || $3 || '%' OR
      email      ILIKE '%' || $3 || '%'
  ))
ORDER BY created_at DESC
LIMIT $4 OFFSET $5
