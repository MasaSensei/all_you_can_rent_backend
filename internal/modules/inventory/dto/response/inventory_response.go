package response

import "time"

// Category is the API-facing shape of entity.Category.
type Category struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenant_id"`
	ParentID    *string     `json:"parent_id,omitempty"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description *string     `json:"description,omitempty"`
	Icon        *string     `json:"icon,omitempty"`
	SortOrder   int         `json:"sort_order"`
	Status      string      `json:"status"`
	Children    []Category  `json:"children,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TemplateField is the API-facing shape of entity.TemplateField.
type TemplateField struct {
	ID           string  `json:"id"`
	FieldName    string  `json:"field_name"`
	FieldLabel   string  `json:"field_label"`
	FieldType    string  `json:"field_type"`
	IsRequired   bool    `json:"is_required"`
	DefaultValue *string `json:"default_value,omitempty"`
	Options      *string `json:"options,omitempty"`
	SortOrder    int     `json:"sort_order"`
}

// AssetTemplate is the API-facing shape of entity.AssetTemplate.
type AssetTemplate struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id"`
	CategoryID  *string         `json:"category_id,omitempty"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Status      string          `json:"status"`
	Fields      []TemplateField `json:"fields,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// AssetValue is the API-facing shape of entity.AssetValue.
type AssetValue struct {
	TemplateFieldID string  `json:"template_field_id"`
	FieldLabel      string  `json:"field_label"`
	FieldType       string  `json:"field_type"`
	Value           *string `json:"value"`
}

// AssetImage is the API-facing shape of entity.AssetImage.
type AssetImage struct {
	ID        string  `json:"id"`
	URL       string  `json:"url"`
	AltText   *string `json:"alt_text,omitempty"`
	IsPrimary bool    `json:"is_primary"`
	SortOrder int     `json:"sort_order"`
}

// AssetDocument is the API-facing shape of entity.AssetDocument.
type AssetDocument struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	FileURL  string  `json:"file_url"`
	FileType *string `json:"file_type,omitempty"`
	FileSize *int    `json:"file_size,omitempty"`
}

// Asset is the full API-facing shape of an asset with all related data.
type Asset struct {
	ID               string          `json:"id"`
	TenantID         string          `json:"tenant_id"`
	CategoryID       *string         `json:"category_id,omitempty"`
	AssetTemplateID  *string         `json:"asset_template_id,omitempty"`
	Name             string          `json:"name"`
	SKU              *string         `json:"sku,omitempty"`
	SerialNumber     *string         `json:"serial_number,omitempty"`
	Description      *string         `json:"description,omitempty"`
	PurchasePrice    *float64        `json:"purchase_price,omitempty"`
	ReplacementValue *float64        `json:"replacement_value,omitempty"`
	PurchaseDate     *time.Time      `json:"purchase_date,omitempty"`
	Condition        string          `json:"condition"`
	Location         *string         `json:"location,omitempty"`
	Status           string          `json:"status"`
	Values           []AssetValue    `json:"values,omitempty"`
	Images           []AssetImage    `json:"images,omitempty"`
	Documents        []AssetDocument `json:"documents,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// AvailabilityStatus is returned by the availability check endpoint,
// also consumed by the booking module via InventoryService interface.
type AvailabilityStatus struct {
	AssetID   string    `json:"asset_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Available bool      `json:"available"`
}
