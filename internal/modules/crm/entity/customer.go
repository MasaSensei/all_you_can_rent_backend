package entity

import "time"

// Customer mirrors the customers table.
type Customer struct {
	ID                 string     `db:"id"`
	TenantID           string     `db:"tenant_id"`
	FirstName          string     `db:"first_name"`
	LastName           string     `db:"last_name"`
	Email              string     `db:"email"`
	Phone              *string    `db:"phone"`
	CompanyName        *string    `db:"company_name"`
	DateOfBirth        *time.Time `db:"date_of_birth"`
	IDDocumentType     *string    `db:"id_document_type"`
	IDDocumentNumber   *string    `db:"id_document_number"`
	CustomerType       string     `db:"customer_type"` // individual, corporate
	Status             string     `db:"status"`
	CreatedBy          *string    `db:"created_by"`
	UpdatedBy          *string    `db:"updated_by"`
	DeletedBy          *string    `db:"deleted_by"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
	Version            int        `db:"version"`
}

const (
	CustomerTypeIndividual = "individual"
	CustomerTypeCorporate  = "corporate"
)
