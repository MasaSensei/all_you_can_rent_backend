package service

import (
	"context"
	"fmt"
)

// BookingSlot is the calendar-view projection returned to the FE.
type BookingSlot struct {
	ID            string  `db:"id"             json:"id"`
	BookingNumber string  `db:"booking_number" json:"booking_number"`
	CustomerName  string  `db:"customer_name"  json:"customer_name"`
	AssetName     string  `db:"asset_name"     json:"asset_name"`
	AssetID       string  `db:"asset_id"       json:"asset_id"`
	CategoryID    *string `db:"category_id"    json:"category_id,omitempty"`
	StartDate     string  `db:"start_date"     json:"start_date"`
	EndDate       string  `db:"end_date"       json:"end_date"`
	Status        string  `db:"status"         json:"status"`
	TotalAmount   float64 `db:"total_amount"   json:"total_amount"`
}

const calendarSQL = `
SELECT
  b.id,
  b.booking_number,
  COALESCE(
    NULLIF(TRIM(COALESCE(c.first_name,'') || ' ' || COALESCE(c.last_name,'')), ''),
    c.company_name,
    'Unknown'
  )                 AS customer_name,
  a.name            AS asset_name,
  a.id              AS asset_id,
  a.category_id,
  b.start_date::text,
  b.end_date::text,
  b.status,
  b.total_amount
FROM bookings  b
JOIN assets    a ON a.id = b.asset_id
JOIN customers c ON c.id = b.customer_id
WHERE b.tenant_id   = $1
  AND b.deleted_at  IS NULL
  AND b.status     != 'cancelled'
  AND b.end_date   >= $2::date
  AND b.start_date <= $3::date
  AND ($4 = '' OR a.category_id::text = $4)
  AND ($5 = '' OR b.status            = $5)
ORDER BY b.start_date ASC
`

// CalendarView returns bookings that overlap [startDate, endDate].
// Implemented on *bookingService so h.booking.CalendarView() works in the handler.
func (s *bookingService) CalendarView(
	ctx context.Context,
	tenantID, startDate, endDate, categoryID, status string,
) ([]BookingSlot, error) {
	var slots []BookingSlot
	if err := s.db.SelectContext(ctx, &slots, calendarSQL,
		tenantID, startDate, endDate, categoryID, status,
	); err != nil {
		return nil, fmt.Errorf("bookingService.CalendarView: %w", err)
	}
	return slots, nil
}
