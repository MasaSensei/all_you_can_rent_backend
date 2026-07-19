package service

import (
	"context"
	"time"

	"rentos/internal/modules/booking/dto/request"
	"rentos/internal/modules/booking/dto/response"
)

// BookingService manages the full booking lifecycle.
type BookingService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateBooking) (*response.Booking, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Booking, error)
	List(ctx context.Context, tenantID string, filter request.ListBookingsFilter) ([]response.Booking, error)
	Confirm(ctx context.Context, id, tenantID, actorID string) (*response.Booking, error)
	Cancel(ctx context.Context, id, tenantID, actorID string, req request.CancelBooking) (*response.Booking, error)
}

// BookingItemService manages extensions and returns for individual items.
type BookingItemService interface {
	Extend(ctx context.Context, itemID, tenantID, actorID string, req request.ExtendBookingItem) (*response.BookingExtension, error)
	Return(ctx context.Context, itemID, tenantID, actorID string, req request.ReturnBookingItem) (*response.BookingReturn, error)
}

// ---- Cross-module dependency interfaces ----
// The booking service depends on these interfaces, never on concrete
// implementations from other modules. This keeps the dependency arrow
// pointing inward and makes each module independently testable.

// InventoryChecker is the subset of inventory.AssetService the booking
// module needs for availability validation.
type InventoryChecker interface {
	CheckAvailability(ctx context.Context, assetID string, start, end time.Time) (bool, error)
}

// PricingQuoter is the subset of pricing.PricingService the booking
// module needs to compute line totals. Implemented in Phase 5.
// Until then, the module boots with a PassthroughPricer that returns
// zero discounts so bookings can be created end-to-end.
type PricingQuoter interface {
	// QuoteItem returns the unit price for one asset over the given range.
	QuoteItem(ctx context.Context, tenantID, assetID string, start, end time.Time, qty int) (unitPrice float64, err error)
	// ValidateCoupon validates a coupon code and returns the discount
	// amount. Returns (0, nil) if code is empty.
	ValidateCoupon(ctx context.Context, tenantID, couponCode string, subtotal float64) (discount float64, couponID *string, err error)
}
