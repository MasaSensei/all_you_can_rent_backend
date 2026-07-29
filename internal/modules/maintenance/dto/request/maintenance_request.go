package request

import "time"

// ScheduleMaintenance creates a new maintenance record.
type ScheduleMaintenance struct {
	AssetID         string     `json:"asset_id" validate:"required,uuid"`
	MaintenanceType string     `json:"maintenance_type" validate:"required,oneof=routine repair cleaning inspection"`
	Description     *string    `json:"description" validate:"omitempty,max=1000"`
	Cost            float64    `json:"cost" validate:"min=0"`
	ScheduledDate   *time.Time `json:"scheduled_date"`
	PerformedBy     *string    `json:"performed_by" validate:"omitempty,max=150"`
}

// UpdateMaintenanceStatus updates the status of a maintenance record.
type UpdateMaintenanceStatus struct {
	MaintenanceStatus string     `json:"maintenance_status" validate:"required,oneof=scheduled in_progress completed cancelled"`
	CompletedDate     *time.Time `json:"completed_date"`
	Cost              *float64   `json:"cost" validate:"omitempty,min=0"`
	PerformedBy       *string    `json:"performed_by" validate:"omitempty,max=150"`
}

// ListMaintenanceFilter holds whitelisted query-param filters.
type ListMaintenanceFilter struct {
	AssetID           *string
	MaintenanceStatus *string
	Page              int
	PerPage           int
}

// CreateInspection records a new asset inspection.
type CreateInspection struct {
	AssetID        string  `json:"asset_id" validate:"required,uuid"`
	BookingItemID  *string `json:"booking_item_id" validate:"omitempty,uuid"`
	InspectionType string  `json:"inspection_type" validate:"required,oneof=pre_rental post_rental routine"`
	InspectorName  *string `json:"inspector_name" validate:"omitempty,max=150"`
	Findings       *string `json:"findings" validate:"omitempty,max=2000"`
	Result         string  `json:"result" validate:"required,oneof=pass fail conditional"`
}

// ListInspectionsFilter holds whitelisted query-param filters.
type ListInspectionsFilter struct {
	AssetID        *string
	InspectionType *string
	Result         *string
	Page           int
	PerPage        int
}

// CreateDamageReport opens a new damage report.
type CreateDamageReport struct {
	AssetID       string  `json:"asset_id" validate:"required,uuid"`
	BookingID     *string `json:"booking_id" validate:"omitempty,uuid"`
	InspectionID  *string `json:"inspection_id" validate:"omitempty,uuid"`
	Description   string  `json:"description" validate:"required,max=2000"`
	Severity      string  `json:"severity" validate:"required,oneof=minor moderate severe total_loss"`
	RepairCost    float64 `json:"repair_cost" validate:"min=0"`
	ChargedAmount float64 `json:"charged_amount" validate:"min=0"`
}

// UpdateDamageReportStatus updates report status and repair cost.
type UpdateDamageReportStatus struct {
	ReportStatus  string  `json:"report_status" validate:"required,oneof=open in_repair resolved closed"`
	RepairCost    *float64 `json:"repair_cost" validate:"omitempty,min=0"`
	ChargedAmount *float64 `json:"charged_amount" validate:"omitempty,min=0"`
}

// ListDamageReportsFilter holds whitelisted query-param filters.
type ListDamageReportsFilter struct {
	AssetID      *string
	ReportStatus *string
	Severity     *string
	Page         int
	PerPage      int
}
