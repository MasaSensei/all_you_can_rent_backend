package response

import "time"

// BookingItem is the API-facing shape of entity.BookingItem.
type BookingItem struct {
	ID        string    `json:"id"`
	AssetID   string    `json:"asset_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	LineTotal float64   `json:"line_total"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

// Booking is the full API-facing shape of a booking.
type Booking struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	CustomerID    string        `json:"customer_id"`
	CouponID      *string       `json:"coupon_id,omitempty"`
	BookingNumber string        `json:"booking_number"`
	StartDate     time.Time     `json:"start_date"`
	EndDate       time.Time     `json:"end_date"`
	Subtotal      float64       `json:"subtotal"`
	TaxTotal      float64       `json:"tax_total"`
	DiscountTotal float64       `json:"discount_total"`
	TotalAmount   float64       `json:"total_amount"`
	BookingStatus string        `json:"booking_status"`
	PaymentStatus string        `json:"payment_status"`
	Notes         *string       `json:"notes,omitempty"`
	Items         []BookingItem `json:"items,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// BookingExtension is the API-facing shape of entity.BookingExtension.
type BookingExtension struct {
	ID             string    `json:"id"`
	BookingID      string    `json:"booking_id"`
	BookingItemID  string    `json:"booking_item_id"`
	OldEndDate     time.Time `json:"old_end_date"`
	NewEndDate     time.Time `json:"new_end_date"`
	AdditionalCost float64   `json:"additional_cost"`
	Reason         *string   `json:"reason,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// BookingReturn is the API-facing shape of entity.BookingReturn.
type BookingReturn struct {
	ID                string    `json:"id"`
	BookingID         string    `json:"booking_id"`
	BookingItemID     string    `json:"booking_item_id"`
	ReturnedAt        time.Time `json:"returned_at"`
	ConditionOnReturn string    `json:"condition_on_return"`
	LateFee           float64   `json:"late_fee"`
	DamageFee         float64   `json:"damage_fee"`
	Notes             *string   `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
