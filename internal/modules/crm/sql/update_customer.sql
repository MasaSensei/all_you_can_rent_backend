UPDATE customers
SET first_name        = COALESCE($3, first_name),
    last_name         = COALESCE($4, last_name),
    phone             = COALESCE($5, phone),
    company_name      = COALESCE($6, company_name),
    date_of_birth     = COALESCE($7, date_of_birth),
    id_document_type  = COALESCE($8, id_document_type),
    id_document_number = COALESCE($9, id_document_number),
    updated_by        = $2,
    updated_at        = now(),
    version           = version + 1
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
