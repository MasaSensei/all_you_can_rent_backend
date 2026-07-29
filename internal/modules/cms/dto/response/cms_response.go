package response

import "time"

// Website is the API-facing shape of entity.Website.
type Website struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Domain    string    `json:"domain"`
	Title     string    `json:"title"`
	Theme     *string   `json:"theme,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Page is the API-facing shape of entity.Page.
type Page struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	WebsiteID string    `json:"website_id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   *string   `json:"content,omitempty"`
	Template  *string   `json:"template,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MenuItem is the API-facing shape of entity.MenuItem — includes children
// for tree rendering.
type MenuItem struct {
	ID        string      `json:"id"`
	ParentID  *string     `json:"parent_id,omitempty"`
	Label     string      `json:"label"`
	URL       *string     `json:"url,omitempty"`
	SortOrder int         `json:"sort_order"`
	Children  []MenuItem  `json:"children,omitempty"`
}

// Menu is the API-facing shape of entity.Menu with its item tree.
type Menu struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	WebsiteID string     `json:"website_id"`
	Name      string     `json:"name"`
	Location  *string    `json:"location,omitempty"`
	Status    string     `json:"status"`
	Items     []MenuItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// BlogCategory is the API-facing shape of entity.BlogCategory.
type BlogCategory struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Blog is the API-facing shape of entity.Blog.
type Blog struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	WebsiteID      string     `json:"website_id"`
	AuthorID       *string    `json:"author_id,omitempty"`
	BlogCategoryID *string    `json:"blog_category_id,omitempty"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	Content        *string    `json:"content,omitempty"`
	FeaturedImage  *string    `json:"featured_image,omitempty"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// SEOMeta is the API-facing shape of entity.SEOMeta.
type SEOMeta struct {
	ID              string  `json:"id"`
	EntityType      string  `json:"entity_type"`
	EntityID        string  `json:"entity_id"`
	MetaTitle       *string `json:"meta_title,omitempty"`
	MetaDescription *string `json:"meta_description,omitempty"`
	MetaKeywords    *string `json:"meta_keywords,omitempty"`
	OGImage         *string `json:"og_image,omitempty"`
	CanonicalURL    *string `json:"canonical_url,omitempty"`
}
