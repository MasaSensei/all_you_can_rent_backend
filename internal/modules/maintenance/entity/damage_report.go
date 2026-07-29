package entity

import "time"

// DamageReport mirrors the damage_reports table.
type DamageReport struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	AssetID      string     `db:"asset_id"`
	BookingID    *string    `db:"booking_id"`
	InspectionID *string    `db:"inspection_id"`
	Description  string     `db:"description"`
	Severity     string     `db:"severity"` // minor, moderate, severe, total_loss
	RepairCost   float64    `db:"repair_cost"`
	ChargedAmount float64   `db:"charged_amount"`
	ReportStatus string     `db:"report_status"` // open, in_repair, resolved, closed
	Status       string     `db:"status"`
	CreatedBy    *string    `db:"created_by"`
	UpdatedBy    *string    `db:"updated_by"`
	DeletedBy    *string    `db:"deleted_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

const (
	SeverityMinor     = "minor"
	SeverityModerate  = "moderate"
	SeveritySevere    = "severe"
	SeverityTotalLoss = "total_loss"

	ReportStatusOpen     = "open"
	ReportStatusInRepair = "in_repair"
	ReportStatusResolved = "resolved"
	ReportStatusClosed   = "closed"
)
