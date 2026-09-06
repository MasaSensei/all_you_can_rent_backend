package response

import "time"

// SuperAdminAuth is returned after successful super admin login.
type SuperAdminAuth struct {
	AccessToken string      `json:"access_token"`
	Admin       SuperAdminInfo `json:"admin"`
}

type SuperAdminInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// TenantListItem is the compact view used in the admin tenant list.
type TenantListItem struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Slug               string     `json:"slug"`
	Email              string     `json:"email"`
	Phone              *string    `json:"phone,omitempty"`
	SubscriptionStatus string     `json:"subscription_status"`
	TrialEndsAt        *time.Time `json:"trial_ends_at,omitempty"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	// Enriched
	ActivePlan      *string    `json:"active_plan,omitempty"`
	PlanEndsAt      *time.Time `json:"plan_ends_at,omitempty"`
	TotalUsers      int        `json:"total_users"`
	TotalAssets     int        `json:"total_assets"`
	TotalBookings   int        `json:"total_bookings"`
}

// TenantDetail is the full view for a single tenant.
type TenantDetail struct {
	TenantListItem
	Timezone string              `json:"timezone"`
	Locale   string              `json:"locale"`
	Currency string              `json:"currency"`
	Subscription *SubscriptionInfo `json:"subscription,omitempty"`
}

type SubscriptionInfo struct {
	ID           string    `json:"id"`
	PlanName     string    `json:"plan_name"`
	PlanSlug     string    `json:"plan_slug"`
	BillingCycle string    `json:"billing_cycle"`
	PricePaid    float64   `json:"price_paid"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	AutoRenew    bool      `json:"auto_renew"`
	Status       string    `json:"status"`
}

// SubscriptionPlan is the public-facing plan info.
type SubscriptionPlan struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	Description  *string  `json:"description,omitempty"`
	PriceMonthly float64  `json:"price_monthly"`
	PriceYearly  float64  `json:"price_yearly"`
	MaxAssets    *int     `json:"max_assets,omitempty"`
	MaxUsers     *int     `json:"max_users,omitempty"`
	Features     []string `json:"features"`
}

// RegisterTenantResult is returned after successful registration.
type RegisterTenantResult struct {
	TenantID   string `json:"tenant_id"`
	TenantSlug string `json:"tenant_slug"`
	AdminEmail string `json:"admin_email"`
	Plan       string `json:"plan"`
	Message    string `json:"message"`
}

// PlatformStats is the super admin overview.
type PlatformStats struct {
	TotalTenants       int     `json:"total_tenants"`
	ActiveTenants      int     `json:"active_tenants"`
	TrialTenants       int     `json:"trial_tenants"`
	TotalMRR           float64 `json:"total_mrr"`
	NewTenantsThisMonth int    `json:"new_tenants_this_month"`
}
