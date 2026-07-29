package entity

import "time"

// Inspection mirrors the inspections table.
type Inspection struct {
	ID             string     `db:"id"`
	TenantID       string     `db:"tenant_id"`
	AssetID        string     `db:"asset_id"`
	BookingItemID  *string    `db:"booking_item_id"`
	InspectionType string     `db:"inspection_type"` // pre_rental, post_rental, routine
	InspectedAt    time.Time  `db:"inspected_at"`
	InspectorName  *string    `db:"inspector_name"`
	Findings       *string    `db:"findings"`
	Result         string     `db:"result"` // pass, fail, conditional
	Status         string     `db:"status"`
	CreatedBy      *string    `db:"created_by"`
	UpdatedBy      *string    `db:"updated_by"`
	DeletedBy      *string    `db:"deleted_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
	Version        int        `db:"version"`
}

const (
	InspectionTypePreRental  = "pre_rental"
	InspectionTypePostRental = "post_rental"
	InspectionTypeRoutine    = "routine"

	InspectionResultPass        = "pass"
	InspectionResultFail        = "fail"
	InspectionResultConditional = "conditional"
)
