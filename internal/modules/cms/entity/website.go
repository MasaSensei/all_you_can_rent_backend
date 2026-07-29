package entity

import "time"

// Website mirrors the websites table.
type Website struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	Domain    string     `db:"domain"`
	Title     string     `db:"title"`
	Theme     *string    `db:"theme"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// Page mirrors the pages table.
type Page struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	WebsiteID string     `db:"website_id"`
	Title     string     `db:"title"`
	Slug      string     `db:"slug"`
	Content   *string    `db:"content"`
	Template  *string    `db:"template"`
	Status    string     `db:"status"` // draft, published
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// Menu mirrors the menus table.
type Menu struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	WebsiteID string     `db:"website_id"`
	Name      string     `db:"name"`
	Location  *string    `db:"location"` // header, footer, sidebar
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// MenuItem mirrors the menu_items table.
type MenuItem struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	MenuID    string     `db:"menu_id"`
	ParentID  *string    `db:"parent_id"`
	Label     string     `db:"label"`
	URL       *string    `db:"url"`
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
