package cms

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/cms/handler"
	"rentos-backend/internal/modules/cms/repository/postgres"
	"rentos-backend/internal/modules/cms/routes"
	"rentos-backend/internal/modules/cms/service"
)

// Module holds the CMS module's wired handler.
type Module struct {
	handler *handler.Handler
}

// New builds the CMS module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	websiteRepo := postgres.NewWebsiteRepository(
		query("create_website.sql"),
		query("find_website_by_id.sql"),
		query("list_websites.sql"),
		query("update_website.sql"),
		query("delete_website.sql"),
	)
	pageRepo := postgres.NewPageRepository(
		query("create_page.sql"),
		query("find_page_by_id.sql"),
		query("find_page_by_slug.sql"),
		query("list_pages.sql"),
		query("update_page.sql"),
		query("delete_page.sql"),
	)
	menuRepo := postgres.NewMenuRepository(
		query("create_menu.sql"),
		query("find_menu_by_id.sql"),
		query("list_menus.sql"),
		query("delete_menu.sql"),
	)
	menuItemRepo := postgres.NewMenuItemRepository(
		query("create_menu_item.sql"),
		query("list_menu_items.sql"),
		query("delete_menu_item.sql"),
	)
	blogCategoryRepo := postgres.NewBlogCategoryRepository(
		query("create_blog_category.sql"),
		query("list_blog_categories.sql"),
		query("delete_blog_category.sql"),
	)
	blogRepo := postgres.NewBlogRepository(
		query("create_blog.sql"),
		query("find_blog_by_id.sql"),
		query("find_blog_by_slug.sql"),
		query("list_blogs.sql"),
		query("update_blog.sql"),
		query("delete_blog.sql"),
	)
	seoRepo := postgres.NewSEOMetaRepository(
		query("upsert_seo_meta.sql"),
		query("find_seo_meta.sql"),
	)

	websiteSvc := service.NewWebsiteService(c.DB, websiteRepo)
	pageSvc := service.NewPageService(c.DB, pageRepo)
	menuSvc := service.NewMenuService(c.DB, menuRepo, menuItemRepo)
	blogSvc := service.NewBlogService(c.DB, blogRepo, blogCategoryRepo)
	seoSvc := service.NewSEOService(c.DB, seoRepo)

	h := handler.New(websiteSvc, pageSvc, menuSvc, blogSvc, seoSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
// Admin routes live under /cms/* (protected by AuthMiddleware in main.go).
// Public read routes live under /public/* (unauthenticated).
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
