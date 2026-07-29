package repository

import (
	"context"
	"errors"
	"time"

	"rentos-backend/internal/modules/booking/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// BookingRepository manages the bookings table.
type BookingRepository interface {
	Create(ctx context.Context, q database.Querier, b *entity.Booking) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Booking, error)
	List(ctx context.Context, q database.Querier, tenantID string, customerID, status *string, limit, offset int) ([]entity.Booking, error)
	UpdateStatus(ctx context.Context, q database.Querier, id, tenantID, status, actorID string) error
	UpdateTotals(ctx context.Context, q database.Querier, b *entity.Booking) error
}

// BookingItemRepository manages the booking_items table.
type BookingItemRepository interface {
	Create(ctx context.Context, q database.Querier, item *entity.BookingItem) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.BookingItem, error)
	ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingItem, error)
	UpdateEndDate(ctx context.Context, q database.Querier, id, tenantID string, newEnd time.Time, actorID string) error
	// CountOverlaps returns how many active booking items conflict with the
	// given asset + date range. Used alongside inventory.CheckAvailability
	// to catch in-flight reservations not yet reflected in asset_availability.
	CountOverlaps(ctx context.Context, q database.Querier, assetID string, start, end time.Time) (int, error)
}

// BookingExtensionRepository manages the booking_extensions table.
type BookingExtensionRepository interface {
	Create(ctx context.Context, q database.Querier, ext *entity.BookingExtension) error
	ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingExtension, error)
}

// BookingReturnRepository manages the booking_returns table.
type BookingReturnRepository interface {
	Create(ctx context.Context, q database.Querier, ret *entity.BookingReturn) error
	ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingReturn, error)
}
