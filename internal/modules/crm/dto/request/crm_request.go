package request

import "time"

// CreateCustomer creates a new customer record.
type CreateCustomer struct {
	FirstName        string     `json:"first_name" validate:"required,max=100"`
	LastName         string     `json:"last_name" validate:"required,max=100"`
	Email            string     `json:"email" validate:"required,email,max=255"`
	Phone            *string    `json:"phone" validate:"omitempty,max=30"`
	CompanyName      *string    `json:"company_name" validate:"omitempty,max=255"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	IDDocumentType   *string    `json:"id_document_type" validate:"omitempty,max=50"`
	IDDocumentNumber *string    `json:"id_document_number" validate:"omitempty,max=100"`
	CustomerType     string     `json:"customer_type" validate:"required,oneof=individual corporate"`
}

// UpdateCustomer updates an existing customer's profile.
type UpdateCustomer struct {
	FirstName        *string    `json:"first_name" validate:"omitempty,max=100"`
	LastName         *string    `json:"last_name" validate:"omitempty,max=100"`
	Phone            *string    `json:"phone" validate:"omitempty,max=30"`
	CompanyName      *string    `json:"company_name" validate:"omitempty,max=255"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	IDDocumentType   *string    `json:"id_document_type" validate:"omitempty,max=50"`
	IDDocumentNumber *string    `json:"id_document_number" validate:"omitempty,max=100"`
}

// ListCustomersFilter holds whitelisted query-param filters.
type ListCustomersFilter struct {
	Search       *string
	CustomerType *string
	Page         int
	PerPage      int
}

// AddAddress adds an address to a customer.
type AddAddress struct {
	AddressType string  `json:"address_type" validate:"required,oneof=billing shipping"`
	Line1       string  `json:"line1" validate:"required,max=255"`
	Line2       *string `json:"line2" validate:"omitempty,max=255"`
	City        string  `json:"city" validate:"required,max=100"`
	State       *string `json:"state" validate:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=20"`
	Country     string  `json:"country" validate:"required,max=100"`
	IsDefault   bool    `json:"is_default"`
}

// CreateMembership creates a membership for a customer.
type CreateMembership struct {
	PlanName  string     `json:"plan_name" validate:"required,max=100"`
	Tier      *string    `json:"tier" validate:"omitempty,max=50"`
	StartDate time.Time  `json:"start_date" validate:"required"`
	EndDate   *time.Time `json:"end_date"`
	Fee       float64    `json:"fee" validate:"min=0"`
}

// CreateLoyaltyProgram creates a new loyalty program.
type CreateLoyaltyProgram struct {
	Name              string  `json:"name" validate:"required,max=150"`
	Description       *string `json:"description" validate:"omitempty,max=500"`
	PointsPerCurrency float64 `json:"points_per_currency" validate:"required,min=0"`
	RedemptionRate    float64 `json:"redemption_rate" validate:"required,min=0"`
}

// EarnPoints records a point-earning transaction.
type EarnPoints struct {
	LoyaltyProgramID string  `json:"loyalty_program_id" validate:"required,uuid"`
	BookingID        *string `json:"booking_id" validate:"omitempty,uuid"`
	Points           int     `json:"points" validate:"required,min=1"`
	Description      *string `json:"description" validate:"omitempty,max=255"`
}

// RedeemPoints records a point-redemption transaction.
type RedeemPoints struct {
	LoyaltyProgramID string  `json:"loyalty_program_id" validate:"required,uuid"`
	Points           int     `json:"points" validate:"required,min=1"`
	Description      *string `json:"description" validate:"omitempty,max=255"`
}
