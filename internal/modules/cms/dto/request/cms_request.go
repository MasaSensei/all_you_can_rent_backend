package request

// ---- Website ----

type CreateWebsite struct {
	Domain string  `json:"domain" validate:"required,fqdn,max=255"`
	Title  string  `json:"title" validate:"required,max=255"`
	Theme  *string `json:"theme" validate:"omitempty,max=100"`
}

type UpdateWebsite struct {
	Title *string `json:"title" validate:"omitempty,max=255"`
	Theme *string `json:"theme" validate:"omitempty,max=100"`
}

// ---- Page ----

type CreatePage struct {
	WebsiteID string  `json:"website_id" validate:"required,uuid"`
	Title     string  `json:"title" validate:"required,max=255"`
	Slug      string  `json:"slug" validate:"required,max=255"`
	Content   *string `json:"content"`
	Template  *string `json:"template" validate:"omitempty,max=100"`
}

type UpdatePage struct {
	Title    *string `json:"title" validate:"omitempty,max=255"`
	Content  *string `json:"content"`
	Template *string `json:"template" validate:"omitempty,max=100"`
	Status   *string `json:"status" validate:"omitempty,oneof=draft published"`
}

// ---- Menu ----

type CreateMenu struct {
	WebsiteID string  `json:"website_id" validate:"required,uuid"`
	Name      string  `json:"name" validate:"required,max=150"`
	Location  *string `json:"location" validate:"omitempty,oneof=header footer sidebar"`
}

type AddMenuItem struct {
	ParentID  *string `json:"parent_id" validate:"omitempty,uuid"`
	Label     string  `json:"label" validate:"required,max=150"`
	URL       *string `json:"url" validate:"omitempty,max=500"`
	SortOrder int     `json:"sort_order"`
}

// ---- Blog Category ----

type CreateBlogCategory struct {
	Name string `json:"name" validate:"required,max=150"`
	Slug string `json:"slug" validate:"required,max=150,alphanumunicode"`
}

// ---- Blog ----

type CreateBlog struct {
	WebsiteID      string  `json:"website_id" validate:"required,uuid"`
	BlogCategoryID *string `json:"blog_category_id" validate:"omitempty,uuid"`
	Title          string  `json:"title" validate:"required,max=255"`
	Slug           string  `json:"slug" validate:"required,max=255"`
	Content        *string `json:"content"`
	FeaturedImage  *string `json:"featured_image" validate:"omitempty,url,max=500"`
}

type UpdateBlog struct {
	Title         *string `json:"title" validate:"omitempty,max=255"`
	Content       *string `json:"content"`
	FeaturedImage *string `json:"featured_image" validate:"omitempty,url,max=500"`
	Status        *string `json:"status" validate:"omitempty,oneof=draft published"`
}

// ListBlogsFilter holds whitelisted query-param filters for blog listing.
type ListBlogsFilter struct {
	WebsiteID      *string
	BlogCategoryID *string
	Status         *string
	Page           int
	PerPage        int
}

// ---- SEO ----

type UpsertSEO struct {
	EntityType      string  `json:"entity_type" validate:"required,oneof=page blog asset"`
	EntityID        string  `json:"entity_id" validate:"required,uuid"`
	MetaTitle       *string `json:"meta_title" validate:"omitempty,max=255"`
	MetaDescription *string `json:"meta_description" validate:"omitempty,max=500"`
	MetaKeywords    *string `json:"meta_keywords" validate:"omitempty,max=500"`
	OGImage         *string `json:"og_image" validate:"omitempty,url,max=500"`
	CanonicalURL    *string `json:"canonical_url" validate:"omitempty,url,max=500"`
}
