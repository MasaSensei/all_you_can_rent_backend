package response

import "time"

// MaintenanceRecord is the API-facing shape of entity.MaintenanceRecord.
type MaintenanceRecord struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	AssetID           string     `json:"asset_id"`
	MaintenanceType   string     `json:"maintenance_type"`
	Description       *string    `json:"description,omitempty"`
	Cost              float64    `json:"cost"`
	ScheduledDate     *time.Time `json:"scheduled_date,omitempty"`
	CompletedDate     *time.Time `json:"completed_date,omitempty"`
	PerformedBy       *string    `json:"performed_by,omitempty"`
	MaintenanceStatus string     `json:"maintenance_status"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Inspection is the API-facing shape of entity.Inspection.
type Inspection struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	AssetID        string    `json:"asset_id"`
	BookingItemID  *string   `json:"booking_item_id,omitempty"`
	InspectionType string    `json:"inspection_type"`
	InspectedAt    time.Time `json:"inspected_at"`
	InspectorName  *string   `json:"inspector_name,omitempty"`
	Findings       *string   `json:"findings,omitempty"`
	Result         string    `json:"result"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DamageReport is the API-facing shape of entity.DamageReport.
type DamageReport struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	AssetID       string    `json:"asset_id"`
	BookingID     *string   `json:"booking_id,omitempty"`
	InspectionID  *string   `json:"inspection_id,omitempty"`
	Description   string    `json:"description"`
	Severity      string    `json:"severity"`
	RepairCost    float64   `json:"repair_cost"`
	ChargedAmount float64   `json:"charged_amount"`
	ReportStatus  string    `json:"report_status"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
