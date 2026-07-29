package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/booking/dto/request"
	"rentos-backend/internal/modules/booking/dto/response"
	"rentos-backend/internal/modules/booking/entity"
	"rentos-backend/internal/modules/booking/repository"
	pkgresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ============================================================
// bookingService
// ============================================================

type bookingService struct {
	db        *sqlx.DB
	bookings  repository.BookingRepository
	items     repository.BookingItemRepository
	inventory InventoryChecker
	pricing   PricingQuoter
}

func NewBookingService(
	db *sqlx.DB,
	bookings repository.BookingRepository,
	items repository.BookingItemRepository,
	inventory InventoryChecker,
	pricing PricingQuoter,
) BookingService {
	return &bookingService{
		db: db, bookings: bookings, items: items,
		inventory: inventory, pricing: pricing,
	}
}

func (s *bookingService) Create(ctx context.Context, tenantID, actorID string, req request.CreateBooking) (*response.Booking, error) {
	// ---- 1. Availability + pricing per item ----
	type quotedItem struct {
		input     request.BookingItemInput
		unitPrice float64
		lineTotal float64
	}

	quoted := make([]quotedItem, 0, len(req.Items))
	var subtotal float64

	for _, input := range req.Items {
		// Inventory availability check (asset_availability + active booking items).
		avail, err := s.inventory.CheckAvailability(ctx, input.AssetID, input.StartDate, input.EndDate)
		if err != nil {
			return nil, err
		}
		if !avail {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict,
				fmt.Sprintf("asset %s is not available for the requested period", input.AssetID))
		}

		// In-flight booking conflict check.
		overlaps, err := s.items.CountOverlaps(ctx, s.db, input.AssetID, input.StartDate, input.EndDate)
		if err != nil {
			return nil, err
		}
		if overlaps > 0 {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict,
				fmt.Sprintf("asset %s already has an active booking for the requested period", input.AssetID))
		}

		unitPrice, err := s.pricing.QuoteItem(ctx, tenantID, input.AssetID, input.StartDate, input.EndDate, input.Quantity)
		if err != nil {
			return nil, err
		}
		lineTotal := unitPrice * float64(input.Quantity)
		subtotal += lineTotal
		quoted = append(quoted, quotedItem{input: input, unitPrice: unitPrice, lineTotal: lineTotal})
	}

	// ---- 2. Coupon validation ----
	var couponCode string
	if req.CouponCode != nil {
		couponCode = *req.CouponCode
	}
	discount, couponID, err := s.pricing.ValidateCoupon(ctx, tenantID, couponCode, subtotal)
	if err != nil {
		return nil, err
	}
	total := subtotal - discount

	// ---- 3. Determine overall date range from items ----
	var bookingStart, bookingEnd time.Time
	for i, q := range quoted {
		if i == 0 || q.input.StartDate.Before(bookingStart) {
			bookingStart = q.input.StartDate
		}
		if i == 0 || q.input.EndDate.After(bookingEnd) {
			bookingEnd = q.input.EndDate
		}
	}

	// ---- 4. Persist atomically ----
	bookingID := uuid.NewString()
	b := &entity.Booking{
		ID:            bookingID,
		TenantID:      tenantID,
		CustomerID:    req.CustomerID,
		CouponID:      couponID,
		BookingNumber: generateBookingNumber(),
		StartDate:     bookingStart,
		EndDate:       bookingEnd,
		Subtotal:      subtotal,
		TaxTotal:      0, // applied by finance module on invoice creation
		DiscountTotal: discount,
		TotalAmount:   total,
		Notes:         req.Notes,
		CreatedBy:     &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.bookings.Create(ctx, tx, b); err != nil {
			return err
		}
		for _, q := range quoted {
			item := &entity.BookingItem{
				ID:        uuid.NewString(),
				TenantID:  tenantID,
				BookingID: bookingID,
				AssetID:   q.input.AssetID,
				Quantity:  q.input.Quantity,
				UnitPrice: q.unitPrice,
				LineTotal: q.lineTotal,
				StartDate: q.input.StartDate,
				EndDate:   q.input.EndDate,
				CreatedBy: &actorID,
			}
			if err := s.items.Create(ctx, tx, item); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, bookingID, tenantID)
}

func (s *bookingService) GetByID(ctx context.Context, id, tenantID string) (*response.Booking, error) {
	b, err := s.bookings.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "booking not found")
		}
		return nil, err
	}
	return s.hydrate(ctx, b)
}

func (s *bookingService) List(ctx context.Context, tenantID string, filter request.ListBookingsFilter) ([]response.Booking, error) {
	perPage := filter.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	bookings, err := s.bookings.List(ctx, s.db, tenantID, filter.CustomerID, filter.BookingStatus, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}

	out := make([]response.Booking, 0, len(bookings))
	for _, b := range bookings {
		r, err := s.hydrate(ctx, &b)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, nil
}

func (s *bookingService) Confirm(ctx context.Context, id, tenantID, actorID string) (*response.Booking, error) {
	b, err := s.bookings.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "booking not found")
		}
		return nil, err
	}

	if b.BookingStatus != entity.BookingStatusPending {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict,
			fmt.Sprintf("cannot confirm a booking with status '%s'", b.BookingStatus))
	}

	if err := s.bookings.UpdateStatus(ctx, s.db, id, tenantID, entity.BookingStatusConfirmed, actorID); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *bookingService) Cancel(ctx context.Context, id, tenantID, actorID string, req request.CancelBooking) (*response.Booking, error) {
	b, err := s.bookings.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "booking not found")
		}
		return nil, err
	}

	cancellable := map[string]bool{
		entity.BookingStatusPending:   true,
		entity.BookingStatusConfirmed: true,
	}
	if !cancellable[b.BookingStatus] {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict,
			fmt.Sprintf("cannot cancel a booking with status '%s'", b.BookingStatus))
	}

	if err := s.bookings.UpdateStatus(ctx, s.db, id, tenantID, entity.BookingStatusCancelled, actorID); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

// hydrate loads booking items and assembles the full response.
func (s *bookingService) hydrate(ctx context.Context, b *entity.Booking) (*response.Booking, error) {
	items, err := s.items.ListByBooking(ctx, s.db, b.ID, b.TenantID)
	if err != nil {
		return nil, err
	}

	itemResp := make([]response.BookingItem, 0, len(items))
	for _, item := range items {
		itemResp = append(itemResp, response.BookingItem{
			ID:        item.ID,
			AssetID:   item.AssetID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			LineTotal: item.LineTotal,
			StartDate: item.StartDate,
			EndDate:   item.EndDate,
			Status:    item.Status,
		})
	}

	return &response.Booking{
		ID:            b.ID,
		TenantID:      b.TenantID,
		CustomerID:    b.CustomerID,
		CouponID:      b.CouponID,
		BookingNumber: b.BookingNumber,
		StartDate:     b.StartDate,
		EndDate:       b.EndDate,
		Subtotal:      b.Subtotal,
		TaxTotal:      b.TaxTotal,
		DiscountTotal: b.DiscountTotal,
		TotalAmount:   b.TotalAmount,
		BookingStatus: b.BookingStatus,
		PaymentStatus: b.PaymentStatus,
		Notes:         b.Notes,
		Items:         itemResp,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}, nil
}

// ============================================================
// bookingItemService
// ============================================================

type bookingItemService struct {
	db         *sqlx.DB
	bookings   repository.BookingRepository
	items      repository.BookingItemRepository
	extensions repository.BookingExtensionRepository
	returns    repository.BookingReturnRepository
	inventory  InventoryChecker
	pricing    PricingQuoter
}

func NewBookingItemService(
	db *sqlx.DB,
	bookings repository.BookingRepository,
	items repository.BookingItemRepository,
	extensions repository.BookingExtensionRepository,
	returns repository.BookingReturnRepository,
	inventory InventoryChecker,
	pricing PricingQuoter,
) BookingItemService {
	return &bookingItemService{
		db: db, bookings: bookings, items: items,
		extensions: extensions, returns: returns,
		inventory: inventory, pricing: pricing,
	}
}

func (s *bookingItemService) Extend(ctx context.Context, itemID, tenantID, actorID string, req request.ExtendBookingItem) (*response.BookingExtension, error) {
	item, err := s.items.FindByID(ctx, s.db, itemID, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "booking item not found")
		}
		return nil, err
	}

	if !req.NewEndDate.After(item.EndDate) {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeValidation, "new end date must be after current end date")
	}

	// Check availability for the extended period.
	avail, err := s.inventory.CheckAvailability(ctx, item.AssetID, item.EndDate, req.NewEndDate)
	if err != nil {
		return nil, err
	}
	if !avail {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "asset is not available for the extension period")
	}

	// Price the delta.
	additionalCost, err := s.pricing.QuoteItem(ctx, tenantID, item.AssetID, item.EndDate, req.NewEndDate, item.Quantity)
	if err != nil {
		return nil, err
	}

	ext := &entity.BookingExtension{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		BookingID:      item.BookingID,
		BookingItemID:  item.ID,
		OldEndDate:     item.EndDate,
		NewEndDate:     req.NewEndDate,
		AdditionalCost: additionalCost,
		Reason:         req.Reason,
		CreatedBy:      &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.extensions.Create(ctx, tx, ext); err != nil {
			return err
		}
		return s.items.UpdateEndDate(ctx, tx, itemID, tenantID, req.NewEndDate, actorID)
	}); err != nil {
		return nil, err
	}

	return &response.BookingExtension{
		ID:             ext.ID,
		BookingID:      ext.BookingID,
		BookingItemID:  ext.BookingItemID,
		OldEndDate:     ext.OldEndDate,
		NewEndDate:     ext.NewEndDate,
		AdditionalCost: ext.AdditionalCost,
		Reason:         ext.Reason,
		CreatedAt:      ext.CreatedAt,
	}, nil
}

func (s *bookingItemService) Return(ctx context.Context, itemID, tenantID, actorID string, req request.ReturnBookingItem) (*response.BookingReturn, error) {
	item, err := s.items.FindByID(ctx, s.db, itemID, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "booking item not found")
		}
		return nil, err
	}

	ret := &entity.BookingReturn{
		ID:                uuid.NewString(),
		TenantID:          tenantID,
		BookingID:         item.BookingID,
		BookingItemID:     item.ID,
		ReturnedAt:        time.Now(),
		ConditionOnReturn: req.Condition,
		LateFee:           req.LateFee,
		DamageFee:         req.DamageFee,
		Notes:             req.Notes,
		CreatedBy:         &actorID,
	}

	if err := s.returns.Create(ctx, s.db, ret); err != nil {
		return nil, err
	}

	// If all items are returned, auto-complete the booking.
	// Full check deferred to a background job in Phase 7 (workers).

	return &response.BookingReturn{
		ID:                ret.ID,
		BookingID:         ret.BookingID,
		BookingItemID:     ret.BookingItemID,
		ReturnedAt:        ret.ReturnedAt,
		ConditionOnReturn: ret.ConditionOnReturn,
		LateFee:           ret.LateFee,
		DamageFee:         ret.DamageFee,
		Notes:             ret.Notes,
		CreatedAt:         ret.CreatedAt,
	}, nil
}

// ============================================================
// PassthroughPricer — placeholder until Phase 5 (pricing module)
// Returns a flat zero price so bookings compile end-to-end now.
// Replaced by the real PricingService in module.go once Phase 5 is done.
// ============================================================

type PassthroughPricer struct{}

func (p *PassthroughPricer) QuoteItem(_ context.Context, _, _ string, _, _ time.Time, _ int) (float64, error) {
	return 0, nil
}

func (p *PassthroughPricer) ValidateCoupon(_ context.Context, _, _ string, _ float64) (float64, *string, error) {
	return 0, nil, nil
}

// ============================================================
// helpers
// ============================================================

func generateBookingNumber() string {
	return fmt.Sprintf("BK-%s", uuid.NewString()[:8])
}
