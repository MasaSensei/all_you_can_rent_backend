package entity

import "time"

// Tenant represents a business that subscribes to RentOS.
type Tenant struct {
	ID                 string     `db:"id"`
	Name               string     `db:"name"`
	Slug               string     `db:"slug"`
	Email              string     `db:"email"`
	Phone              *string    `db:"phone"`
	LogoURL            *string    `db:"logo_url"`
	Address            *string    `db:"address"`
	Timezone           string     `db:"timezone"`
	Locale             string     `db:"locale"`
	Currency           string     `db:"currency"`
	SubscriptionStatus string     `db:"subscription_status"`
	TrialEndsAt        *time.Time `db:"trial_ends_at"`
	Status             string     `db:"status"`
	CreatedBy          *string    `db:"created_by"`
	UpdatedBy          *string    `db:"updated_by"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
	Version            int        `db:"version"`
}

const (
	SubscriptionStatusTrial     = "trial"
	SubscriptionStatusActive    = "active"
	SubscriptionStatusPastDue   = "past_due"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusExpired   = "expired"
)

// SubscriptionPlan defines a pricing tier.
type SubscriptionPlan struct {
	ID                  string    `db:"id"`
	Name                string    `db:"name"`
	Slug                string    `db:"slug"`
	Description         *string   `db:"description"`
	PriceMonthly        float64   `db:"price_monthly"`
	PriceYearly         float64   `db:"price_yearly"`
	MaxAssets           *int      `db:"max_assets"`
	MaxUsers            *int      `db:"max_users"`
	MaxBookingsPerMonth *int      `db:"max_bookings_per_month"`
	Features            []byte    `db:"features"` // JSONB
	IsActive            bool      `db:"is_active"`
	Version             int       `db:"version"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

// TenantSubscription tracks the active subscription of a tenant.
type TenantSubscription struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	PlanID       string     `db:"plan_id"`
	BillingCycle string     `db:"billing_cycle"` // monthly | yearly
	PricePaid    float64    `db:"price_paid"`
	StartsAt     time.Time  `db:"starts_at"`
	EndsAt       time.Time  `db:"ends_at"`
	AutoRenew    bool       `db:"auto_renew"`
	Status       string     `db:"status"`
	CancelledAt  *time.Time `db:"cancelled_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

// SuperAdmin is the platform-level administrator.
type SuperAdmin struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	FullName     string     `db:"full_name"`
	IsActive     bool       `db:"is_active"`
	LastLoginAt  *time.Time `db:"last_login_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	Version      int        `db:"version"`
}
