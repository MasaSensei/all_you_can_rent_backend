package request

// SuperAdminLogin authenticates a super admin.
type SuperAdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterTenant is the public tenant registration request.
type RegisterTenant struct {
	// Business info
	BusinessName  string  `json:"business_name" validate:"required,min=2,max=255"`
	BusinessSlug  string  `json:"business_slug" validate:"required,min=2,max=50,alphanum_dash"`
	BusinessPhone *string `json:"business_phone" validate:"omitempty,max=30"`
	// Admin account
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	// Plan
	Plan string `json:"plan" validate:"required,oneof=trial starter professional enterprise"`
}

// UpdateTenantStatus allows super admin to activate/suspend a tenant.
type UpdateTenantStatus struct {
	Status string `json:"status" validate:"required,oneof=active suspended"`
	Reason string `json:"reason" validate:"omitempty,max=500"`
}

// AssignSubscription assigns a plan to a tenant.
type AssignSubscription struct {
	PlanID       string  `json:"plan_id" validate:"required,uuid"`
	BillingCycle string  `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	PricePaid    float64 `json:"price_paid" validate:"min=0"`
	DurationDays int     `json:"duration_days" validate:"required,min=1"`
}

// ListTenantsFilter for filtering tenants.
type ListTenantsFilter struct {
	Search             string
	SubscriptionStatus string
	Status             string
	Page               int
	PerPage            int
}
