package response

import "time"

// CustomerAddress is the API-facing shape of entity.CustomerAddress.
type CustomerAddress struct {
	ID          string  `json:"id"`
	AddressType string  `json:"address_type"`
	Line1       string  `json:"line1"`
	Line2       *string `json:"line2,omitempty"`
	City        string  `json:"city"`
	State       *string `json:"state,omitempty"`
	PostalCode  *string `json:"postal_code,omitempty"`
	Country     string  `json:"country"`
	IsDefault   bool    `json:"is_default"`
}

// Customer is the full API-facing shape of a customer.
type Customer struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenant_id"`
	FirstName        string            `json:"first_name"`
	LastName         string            `json:"last_name"`
	Email            string            `json:"email"`
	Phone            *string           `json:"phone,omitempty"`
	CompanyName      *string           `json:"company_name,omitempty"`
	DateOfBirth      *time.Time        `json:"date_of_birth,omitempty"`
	IDDocumentType   *string           `json:"id_document_type,omitempty"`
	IDDocumentNumber *string           `json:"id_document_number,omitempty"`
	CustomerType     string            `json:"customer_type"`
	Status           string            `json:"status"`
	Addresses        []CustomerAddress `json:"addresses,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// Membership is the API-facing shape of entity.Membership.
type Membership struct {
	ID               string     `json:"id"`
	CustomerID       string     `json:"customer_id"`
	PlanName         string     `json:"plan_name"`
	Tier             *string    `json:"tier,omitempty"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	Fee              float64    `json:"fee"`
	MembershipStatus string     `json:"membership_status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// LoyaltyProgram is the API-facing shape of entity.LoyaltyProgram.
type LoyaltyProgram struct {
	ID                string  `json:"id"`
	TenantID          string  `json:"tenant_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	PointsPerCurrency float64 `json:"points_per_currency"`
	RedemptionRate    float64 `json:"redemption_rate"`
	Status            string  `json:"status"`
}

// LoyaltyTransaction is the API-facing shape of entity.LoyaltyTransaction.
type LoyaltyTransaction struct {
	ID              string    `json:"id"`
	CustomerID      string    `json:"customer_id"`
	BookingID       *string   `json:"booking_id,omitempty"`
	Points          int       `json:"points"`
	TransactionType string    `json:"transaction_type"`
	Description     *string   `json:"description,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// LoyaltyBalance is the customer's current point balance.
type LoyaltyBalance struct {
	CustomerID string `json:"customer_id"`
	Balance    int    `json:"balance"`
}
