package entity

import "time"

// Report mirrors the reports table.
// Generation is async: POST /reports creates a row with status=queued,
// a worker processes it, then updates status=completed + file_url.
type Report struct {
	ID              string     `db:"id"`
	TenantID        string     `db:"tenant_id"`
	Name            string     `db:"name"`
	ReportType      string     `db:"report_type"`
	Parameters      *string    `db:"parameters"` // JSON
	GeneratedFormat string     `db:"generated_format"` // pdf, csv, xlsx
	FileURL         *string    `db:"file_url"`
	Status          string     `db:"status"` // queued, processing, completed, failed
	CreatedBy       *string    `db:"created_by"`
	UpdatedBy       *string    `db:"updated_by"`
	DeletedBy       *string    `db:"deleted_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	Version         int        `db:"version"`
}

const (
	ReportTypeRevenue      = "revenue"
	ReportTypeAssetUsage   = "asset_usage"
	ReportTypeBookings     = "bookings"
	ReportTypeCustomers    = "customers"
	ReportTypeExpenses     = "expenses"

	ReportStatusQueued     = "queued"
	ReportStatusProcessing = "processing"
	ReportStatusCompleted  = "completed"
	ReportStatusFailed     = "failed"

	ReportFormatPDF  = "pdf"
	ReportFormatCSV  = "csv"
	ReportFormatXLSX = "xlsx"
)

// AnalyticsEvent mirrors the analytics_events table.
// High-volume; inserts are batched by the background worker.
type AnalyticsEvent struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	UserID        *string    `db:"user_id"`
	CustomerID    *string    `db:"customer_id"`
	EventName     string     `db:"event_name"`
	EventCategory *string    `db:"event_category"`
	EventData     *string    `db:"event_data"` // JSON
	Source        *string    `db:"source"`
	OccurredAt    time.Time  `db:"occurred_at"`
	Status        string     `db:"status"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	Version       int        `db:"version"`
}
