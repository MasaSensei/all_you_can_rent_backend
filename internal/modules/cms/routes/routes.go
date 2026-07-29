package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/cms/handler"
)

// Register mounts all CMS routes under /api/v1.
// Admin routes (create/update/delete) require AuthMiddleware upstream.
// Public read routes (GetBySlug) are mounted unauthenticated so the
// rendered website can fetch pages and blogs without a JWT.
func Register(router fiber.Router, h *handler.Handler) {
	// ---- Websites ----
	websites := router.Group("/cms/websites")
	websites.Post("/", h.CreateWebsite)
	websites.Get("/", h.ListWebsites)
	websites.Get("/:id", h.GetWebsite)
	websites.Put("/:id", h.UpdateWebsite)
	websites.Delete("/:id", h.DeleteWebsite)

	// ---- Pages ----
	pages := router.Group("/cms/pages")
	pages.Post("/", h.CreatePage)
	pages.Get("/", h.ListPages)
	pages.Get("/:id", h.GetPage)
	pages.Put("/:id", h.UpdatePage)
	pages.Delete("/:id", h.DeletePage)

	// ---- Menus ----
	menus := router.Group("/cms/menus")
	menus.Post("/", h.CreateMenu)
	menus.Get("/", h.ListMenus)
	menus.Get("/:id", h.GetMenu)
	menus.Delete("/:id", h.DeleteMenu)
	menus.Post("/:id/items", h.AddMenuItem)
	menus.Delete("/:id/items/:item_id", h.DeleteMenuItem)

	// ---- Blog Categories ----
	blogCats := router.Group("/cms/blog-categories")
	blogCats.Post("/", h.CreateBlogCategory)
	blogCats.Get("/", h.ListBlogCategories)
	blogCats.Delete("/:id", h.DeleteBlogCategory)

	// ---- Blogs ----
	blogs := router.Group("/cms/blogs")
	blogs.Post("/", h.CreateBlog)
	blogs.Get("/", h.ListBlogs)
	blogs.Get("/:id", h.GetBlog)
	blogs.Put("/:id", h.UpdateBlog)
	blogs.Delete("/:id", h.DeleteBlog)

	// ---- SEO ----
	seo := router.Group("/cms/seo")
	seo.Put("/", h.UpsertSEO)
	seo.Get("/:entity_type/:entity_id", h.GetSEO)

	// ---- Public read-side (unauthenticated) ----
	// Mounted on the same router; the caller (main.go) is responsible for
	// mounting these on a public group without AuthMiddleware.
	public := router.Group("/public")
	public.Get("/pages/:website_id/:slug", h.GetPageBySlug)
	public.Get("/blogs/:website_id/:slug", h.GetBlogBySlug)
}
