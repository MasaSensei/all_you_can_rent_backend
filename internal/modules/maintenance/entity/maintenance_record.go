package entity

import "time"

// MaintenanceRecord mirrors the maintenance_records table.
type MaintenanceRecord struct {
	ID                string     `db:"id"`
	TenantID          string     `db:"tenant_id"`
	AssetID           string     `db:"asset_id"`
	MaintenanceType   string     `db:"maintenance_type"` // routine, repair, cleaning, inspection
	Description       *string    `db:"description"`
	Cost              float64    `db:"cost"`
	ScheduledDate     *time.Time `db:"scheduled_date"`
	CompletedDate     *time.Time `db:"completed_date"`
	PerformedBy       *string    `db:"performed_by"`
	MaintenanceStatus string     `db:"maintenance_status"` // scheduled, in_progress, completed, cancelled
	Status            string     `db:"status"`
	CreatedBy         *string    `db:"created_by"`
	UpdatedBy         *string    `db:"updated_by"`
	DeletedBy         *string    `db:"deleted_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	Version           int        `db:"version"`
}

const (
	MaintenanceTypeRoutine    = "routine"
	MaintenanceTypeRepair     = "repair"
	MaintenanceTypeCleaning   = "cleaning"
	MaintenanceTypeInspection = "inspection"

	MaintenanceStatusScheduled  = "scheduled"
	MaintenanceStatusInProgress = "in_progress"
	MaintenanceStatusCompleted  = "completed"
	MaintenanceStatusCancelled  = "cancelled"
)
