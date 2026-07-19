package request

// CreateCategory creates a new category.
type CreateCategory struct {
	ParentID    *string `json:"parent_id" validate:"omitempty,uuid"`
	Name        string  `json:"name" validate:"required,min=1,max=150"`
	Slug        string  `json:"slug" validate:"required,min=1,max=150,alphanumunicode"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Icon        *string `json:"icon" validate:"omitempty,max=255"`
	SortOrder   int     `json:"sort_order"`
}

// UpdateCategory updates an existing category.
type UpdateCategory struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=150"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Icon        *string `json:"icon" validate:"omitempty,max=255"`
	SortOrder   *int    `json:"sort_order"`
}
