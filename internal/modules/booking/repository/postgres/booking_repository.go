package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"rentos-backend/internal/modules/booking/entity"
	"rentos-backend/internal/modules/booking/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// bookingRepository
// ============================================================

type bookingRepository struct {
	qCreate        string
	qFindByID      string
	qList          string
	qUpdateStatus  string
	qUpdateTotals  string
}

func NewBookingRepository(qCreate, qFindByID, qList, qUpdateStatus, qUpdateTotals string) repository.BookingRepository {
	return &bookingRepository{
		qCreate: qCreate, qFindByID: qFindByID, qList: qList,
		qUpdateStatus: qUpdateStatus, qUpdateTotals: qUpdateTotals,
	}
}

func (r *bookingRepository) Create(ctx context.Context, q database.Querier, b *entity.Booking) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		b.ID, b.TenantID, b.CustomerID, b.CouponID, b.BookingNumber,
		b.StartDate, b.EndDate, b.Subtotal, b.TaxTotal, b.DiscountTotal, b.TotalAmount,
		b.Notes, b.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("bookingRepository.Create: %w", err)
	}
	return nil
}

func (r *bookingRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Booking, error) {
	var b entity.Booking
	if err := q.GetContext(ctx, &b, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("bookingRepository.FindByID: %w", err)
	}
	return &b, nil
}

func (r *bookingRepository) List(ctx context.Context, q database.Querier, tenantID string, customerID, status *string, limit, offset int) ([]entity.Booking, error) {
	var out []entity.Booking
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, customerID, status, limit, offset); err != nil {
		return nil, fmt.Errorf("bookingRepository.List: %w", err)
	}
	return out, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, q database.Querier, id, tenantID, status, actorID string) error {
	res, err := q.ExecContext(ctx, r.qUpdateStatus, id, tenantID, status)
	if err != nil {
		return fmt.Errorf("bookingRepository.UpdateStatus: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *bookingRepository) UpdateTotals(ctx context.Context, q database.Querier, b *entity.Booking) error {
	_, err := q.ExecContext(ctx, r.qUpdateTotals,
		b.ID, b.TenantID, b.Subtotal, b.TaxTotal, b.DiscountTotal, b.TotalAmount, b.CouponID,
	)
	if err != nil {
		return fmt.Errorf("bookingRepository.UpdateTotals: %w", err)
	}
	return nil
}

// ============================================================
// bookingItemRepository
// ============================================================

type bookingItemRepository struct {
	qCreate        string
	qFindByID      string
	qListByBooking string
	qUpdateEndDate string
	qCountOverlaps string
}

func NewBookingItemRepository(qCreate, qFindByID, qListByBooking, qUpdateEndDate, qCountOverlaps string) repository.BookingItemRepository {
	return &bookingItemRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qListByBooking: qListByBooking, qUpdateEndDate: qUpdateEndDate,
		qCountOverlaps: qCountOverlaps,
	}
}

func (r *bookingItemRepository) Create(ctx context.Context, q database.Querier, item *entity.BookingItem) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		item.ID, item.TenantID, item.BookingID, item.AssetID, item.Quantity,
		item.UnitPrice, item.LineTotal, item.StartDate, item.EndDate, item.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("bookingItemRepository.Create: %w", err)
	}
	return nil
}

func (r *bookingItemRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.BookingItem, error) {
	var item entity.BookingItem
	if err := q.GetContext(ctx, &item, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("bookingItemRepository.FindByID: %w", err)
	}
	return &item, nil
}

func (r *bookingItemRepository) ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingItem, error) {
	var out []entity.BookingItem
	if err := q.SelectContext(ctx, &out, r.qListByBooking, bookingID, tenantID); err != nil {
		return nil, fmt.Errorf("bookingItemRepository.ListByBooking: %w", err)
	}
	return out, nil
}

func (r *bookingItemRepository) UpdateEndDate(ctx context.Context, q database.Querier, id, tenantID string, newEnd time.Time, actorID string) error {
	res, err := q.ExecContext(ctx, r.qUpdateEndDate, id, tenantID, newEnd)
	if err != nil {
		return fmt.Errorf("bookingItemRepository.UpdateEndDate: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *bookingItemRepository) CountOverlaps(ctx context.Context, q database.Querier, assetID string, start, end time.Time) (int, error) {
	var count int
	if err := q.GetContext(ctx, &count, r.qCountOverlaps, assetID, start, end); err != nil {
		return 0, fmt.Errorf("bookingItemRepository.CountOverlaps: %w", err)
	}
	return count, nil
}

// ============================================================
// bookingExtensionRepository
// ============================================================

type bookingExtensionRepository struct {
	qCreate        string
	qListByBooking string
}

func NewBookingExtensionRepository(qCreate, qListByBooking string) repository.BookingExtensionRepository {
	return &bookingExtensionRepository{qCreate: qCreate, qListByBooking: qListByBooking}
}

func (r *bookingExtensionRepository) Create(ctx context.Context, q database.Querier, ext *entity.BookingExtension) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		ext.ID, ext.TenantID, ext.BookingID, ext.BookingItemID,
		ext.OldEndDate, ext.NewEndDate, ext.AdditionalCost, ext.Reason, ext.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("bookingExtensionRepository.Create: %w", err)
	}
	return nil
}

func (r *bookingExtensionRepository) ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingExtension, error) {
	var out []entity.BookingExtension
	if err := q.SelectContext(ctx, &out, r.qListByBooking, bookingID, tenantID); err != nil {
		return nil, fmt.Errorf("bookingExtensionRepository.ListByBooking: %w", err)
	}
	return out, nil
}

// ============================================================
// bookingReturnRepository
// ============================================================

type bookingReturnRepository struct {
	qCreate        string
	qListByBooking string
}

func NewBookingReturnRepository(qCreate, qListByBooking string) repository.BookingReturnRepository {
	return &bookingReturnRepository{qCreate: qCreate, qListByBooking: qListByBooking}
}

func (r *bookingReturnRepository) Create(ctx context.Context, q database.Querier, ret *entity.BookingReturn) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		ret.ID, ret.TenantID, ret.BookingID, ret.BookingItemID,
		ret.ConditionOnReturn, ret.LateFee, ret.DamageFee, ret.Notes, ret.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("bookingReturnRepository.Create: %w", err)
	}
	return nil
}

func (r *bookingReturnRepository) ListByBooking(ctx context.Context, q database.Querier, bookingID, tenantID string) ([]entity.BookingReturn, error) {
	var out []entity.BookingReturn
	if err := q.SelectContext(ctx, &out, r.qListByBooking, bookingID, tenantID); err != nil {
		return nil, fmt.Errorf("bookingReturnRepository.ListByBooking: %w", err)
	}
	return out, nil
}
