INSERT INTO customers (
    id, tenant_id, first_name, last_name, email, phone, company_name,
    date_of_birth, id_document_type, id_document_number,
    customer_type, status, created_by, updated_by, created_at, updated_at, version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10,
    $11, 'active', $12, $12, now(), now(), 1
)
