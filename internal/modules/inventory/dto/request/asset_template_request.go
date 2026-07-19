package request

// CreateAssetTemplate creates a new asset template with its fields.
type CreateAssetTemplate struct {
	CategoryID  *string              `json:"category_id" validate:"omitempty,uuid"`
	Name        string               `json:"name" validate:"required,min=1,max=150"`
	Description *string              `json:"description" validate:"omitempty,max=500"`
	Fields      []CreateTemplateField `json:"fields"`
}

// UpdateAssetTemplate updates an existing template's metadata.
type UpdateAssetTemplate struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=150"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// CreateTemplateField adds a field to a template.
type CreateTemplateField struct {
	FieldName    string  `json:"field_name" validate:"required,max=100"`
	FieldLabel   string  `json:"field_label" validate:"required,max=150"`
	FieldType    string  `json:"field_type" validate:"required,oneof=string number boolean date enum"`
	IsRequired   bool    `json:"is_required"`
	DefaultValue *string `json:"default_value"`
	Options      *string `json:"options"` // JSON array for enum type
	SortOrder    int     `json:"sort_order"`
}
