package request

import "time"

// CreateAsset creates a new rentable asset.
type CreateAsset struct {
	CategoryID       *string            `json:"category_id" validate:"omitempty,uuid"`
	AssetTemplateID  *string            `json:"asset_template_id" validate:"omitempty,uuid"`
	Name             string             `json:"name" validate:"required,min=1,max=255"`
	SKU              *string            `json:"sku" validate:"omitempty,max=100"`
	SerialNumber     *string            `json:"serial_number" validate:"omitempty,max=100"`
	Description      *string            `json:"description"`
	PurchasePrice    *float64           `json:"purchase_price" validate:"omitempty,min=0"`
	ReplacementValue *float64           `json:"replacement_value" validate:"omitempty,min=0"`
	PurchaseDate     *time.Time         `json:"purchase_date"`
	Condition        string             `json:"condition" validate:"required,oneof=new good fair poor"`
	Location         *string            `json:"location" validate:"omitempty,max=255"`
	Values           []AssetValueInput  `json:"values"`
}

// UpdateAsset updates an existing asset's metadata.
type UpdateAsset struct {
	Name             *string  `json:"name" validate:"omitempty,min=1,max=255"`
	Description      *string  `json:"description"`
	Condition        *string  `json:"condition" validate:"omitempty,oneof=new good fair poor"`
	Location         *string  `json:"location" validate:"omitempty,max=255"`
	PurchasePrice    *float64 `json:"purchase_price" validate:"omitempty,min=0"`
	ReplacementValue *float64 `json:"replacement_value" validate:"omitempty,min=0"`
}

// AssetValueInput carries a single EAV field value.
type AssetValueInput struct {
	TemplateFieldID string  `json:"template_field_id" validate:"required,uuid"`
	Value           *string `json:"value"`
}

// ListAssetsFilter holds whitelisted query-param filters.
type ListAssetsFilter struct {
	CategoryID *string
	Status     *string
	Search     *string
	Page       int
	PerPage    int
}

// BlockAvailability blocks a date range on an asset.
type BlockAvailability struct {
	StartDate        time.Time `json:"start_date" validate:"required"`
	EndDate          time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	AvailabilityType string    `json:"availability_type" validate:"required,oneof=blocked maintenance"`
	Reason           *string   `json:"reason" validate:"omitempty,max=500"`
}

// CheckAvailability is used by the booking module to query a date range.
type CheckAvailability struct {
	AssetID   string    `json:"asset_id" validate:"required,uuid"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}
