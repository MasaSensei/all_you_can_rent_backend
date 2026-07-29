package entity

import "time"

// BlogCategory mirrors the blog_categories table.
type BlogCategory struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	Name      string     `db:"name"`
	Slug      string     `db:"slug"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

// Blog mirrors the blogs table.
type Blog struct {
	ID             string     `db:"id"`
	TenantID       string     `db:"tenant_id"`
	WebsiteID      string     `db:"website_id"`
	AuthorID       *string    `db:"author_id"`
	BlogCategoryID *string    `db:"blog_category_id"`
	Title          string     `db:"title"`
	Slug           string     `db:"slug"`
	Content        *string    `db:"content"`
	FeaturedImage  *string    `db:"featured_image"`
	PublishedAt    *time.Time `db:"published_at"`
	Status         string     `db:"status"` // draft, published
	CreatedBy      *string    `db:"created_by"`
	UpdatedBy      *string    `db:"updated_by"`
	DeletedBy      *string    `db:"deleted_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
	Version        int        `db:"version"`
}

// SEOMeta mirrors the seo_meta table (polymorphic via entity_type+entity_id).
type SEOMeta struct {
	ID              string     `db:"id"`
	TenantID        string     `db:"tenant_id"`
	EntityType      string     `db:"entity_type"` // page, blog, asset
	EntityID        string     `db:"entity_id"`
	MetaTitle       *string    `db:"meta_title"`
	MetaDescription *string    `db:"meta_description"`
	MetaKeywords    *string    `db:"meta_keywords"`
	OGImage         *string    `db:"og_image"`
	CanonicalURL    *string    `db:"canonical_url"`
	Status          string     `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	Version         int        `db:"version"`
}

const (
	ContentStatusDraft     = "draft"
	ContentStatusPublished = "published"

	SEOEntityTypePage  = "page"
	SEOEntityTypeBlog  = "blog"
	SEOEntityTypeAsset = "asset"
)
