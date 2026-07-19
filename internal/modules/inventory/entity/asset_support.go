package entity

import "time"

// AssetValue mirrors the asset_values table (EAV pattern).
type AssetValue struct {
	ID              string     `db:"id"`
	TenantID        string     `db:"tenant_id"`
	AssetID         string     `db:"asset_id"`
	TemplateFieldID string     `db:"template_field_id"`
	Value           *string    `db:"value"`
	Status          string     `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	Version         int        `db:"version"`
}

// AssetImage mirrors the asset_images table.
type AssetImage struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	AssetID   string     `db:"asset_id"`
	URL       string     `db:"url"`
	AltText   *string    `db:"alt_text"`
	IsPrimary bool       `db:"is_primary"`
	SortOrder int        `db:"sort_order"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// AssetDocument mirrors the asset_documents table.
type AssetDocument struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	AssetID   string     `db:"asset_id"`
	Title     string     `db:"title"`
	FileURL   string     `db:"file_url"`
	FileType  *string    `db:"file_type"`
	FileSize  *int       `db:"file_size"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// AssetAvailability mirrors the asset_availability table.
type AssetAvailability struct {
	ID               string     `db:"id"`
	TenantID         string     `db:"tenant_id"`
	AssetID          string     `db:"asset_id"`
	StartDate        time.Time  `db:"start_date"`
	EndDate          time.Time  `db:"end_date"`
	AvailabilityType string     `db:"availability_type"`
	Reason           *string    `db:"reason"`
	Status           string     `db:"status"`
	CreatedBy        *string    `db:"created_by"`
	UpdatedBy        *string    `db:"updated_by"`
	DeletedBy        *string    `db:"deleted_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Version          int        `db:"version"`
}
