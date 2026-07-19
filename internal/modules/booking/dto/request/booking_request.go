package request

import "time"

// BookingItemInput is a single line in a create booking request.
type BookingItemInput struct {
	AssetID   string    `json:"asset_id" validate:"required,uuid"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

// CreateBooking is the input for starting a new booking.
type CreateBooking struct {
	CustomerID  string             `json:"customer_id" validate:"required,uuid"`
	Items       []BookingItemInput `json:"items" validate:"required,min=1,dive"`
	CouponCode  *string            `json:"coupon_code" validate:"omitempty,max=50"`
	Notes       *string            `json:"notes" validate:"omitempty,max=1000"`
}

// ConfirmBooking transitions a booking from pending → confirmed.
type ConfirmBooking struct {
	BookingID string `json:"booking_id" validate:"required,uuid"`
}

// CancelBooking transitions a booking to cancelled.
type CancelBooking struct {
	Reason *string `json:"reason" validate:"omitempty,max=500"`
}

// ExtendBookingItem requests a new end date for a single booking item.
type ExtendBookingItem struct {
	NewEndDate time.Time `json:"new_end_date" validate:"required"`
	Reason     *string   `json:"reason" validate:"omitempty,max=500"`
}

// ReturnBookingItem records the return of a single booking item.
type ReturnBookingItem struct {
	Condition string  `json:"condition" validate:"required,oneof=new good fair poor damaged"`
	LateFee   float64 `json:"late_fee" validate:"min=0"`
	DamageFee float64 `json:"damage_fee" validate:"min=0"`
	Notes     *string `json:"notes" validate:"omitempty,max=1000"`
}

// ListBookingsFilter holds whitelisted query-param filters.
type ListBookingsFilter struct {
	CustomerID    *string
	BookingStatus *string
	Page          int
	PerPage       int
}
