SELECT id, tenant_id, first_name, last_name, email, phone, company_name,
       date_of_birth, id_document_type, id_document_number,
       customer_type, status, created_by, updated_by, deleted_by,
       created_at, updated_at, deleted_at, version
FROM customers
WHERE email = $1 AND tenant_id = $2 AND deleted_at IS NULL
