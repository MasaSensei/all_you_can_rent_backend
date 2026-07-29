package handler

import (
	"github.com/gofiber/fiber/v2"

	cmsreq "rentos-backend/internal/modules/cms/dto/request"
	"rentos-backend/internal/modules/cms/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all CMS HTTP handlers.
type Handler struct {
	websites service.WebsiteService
	pages    service.PageService
	menus    service.MenuService
	blogs    service.BlogService
	seo      service.SEOService
	validate *validator.Validate
}

func New(
	websites service.WebsiteService,
	pages service.PageService,
	menus service.MenuService,
	blogs service.BlogService,
	seo service.SEOService,
	v *validator.Validate,
) *Handler {
	return &Handler{
		websites: websites, pages: pages, menus: menus,
		blogs: blogs, seo: seo, validate: v,
	}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals(ctxKeyTenantID).(string); ok {
		return id
	}
	return c.Get("X-Tenant-ID")
}

func userID(c *fiber.Ctx) string {
	id, _ := c.Locals(ctxKeyUserID).(string)
	return id
}

// ---- Websites ----

func (h *Handler) CreateWebsite(c *fiber.Ctx) error {
	var req cmsreq.CreateWebsite
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	w, err := h.websites.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, w)
}

func (h *Handler) GetWebsite(c *fiber.Ctx) error {
	w, err := h.websites.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, w)
}

func (h *Handler) ListWebsites(c *fiber.Ctx) error {
	websites, err := h.websites.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, websites)
}

func (h *Handler) UpdateWebsite(c *fiber.Ctx) error {
	var req cmsreq.UpdateWebsite
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	w, err := h.websites.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, w)
}

func (h *Handler) DeleteWebsite(c *fiber.Ctx) error {
	if err := h.websites.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Pages ----

func (h *Handler) CreatePage(c *fiber.Ctx) error {
	var req cmsreq.CreatePage
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	p, err := h.pages.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, p)
}

func (h *Handler) GetPage(c *fiber.Ctx) error {
	p, err := h.pages.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, p)
}

func (h *Handler) GetPageBySlug(c *fiber.Ctx) error {
	p, err := h.pages.GetBySlug(c.Context(), c.Params("website_id"), c.Params("slug"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, p)
}

func (h *Handler) ListPages(c *fiber.Ctx) error {
	pages, err := h.pages.List(c.Context(), c.Query("website_id"), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, pages)
}

func (h *Handler) UpdatePage(c *fiber.Ctx) error {
	var req cmsreq.UpdatePage
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	p, err := h.pages.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, p)
}

func (h *Handler) DeletePage(c *fiber.Ctx) error {
	if err := h.pages.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Menus ----

func (h *Handler) CreateMenu(c *fiber.Ctx) error {
	var req cmsreq.CreateMenu
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	m, err := h.menus.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, m)
}

func (h *Handler) GetMenu(c *fiber.Ctx) error {
	m, err := h.menus.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, m)
}

func (h *Handler) ListMenus(c *fiber.Ctx) error {
	menus, err := h.menus.List(c.Context(), c.Query("website_id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, menus)
}

func (h *Handler) DeleteMenu(c *fiber.Ctx) error {
	if err := h.menus.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) AddMenuItem(c *fiber.Ctx) error {
	var req cmsreq.AddMenuItem
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	item, err := h.menus.AddItem(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, item)
}

func (h *Handler) DeleteMenuItem(c *fiber.Ctx) error {
	if err := h.menus.DeleteItem(c.Context(), c.Params("item_id"), c.Params("id"), tenantID(c), userID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Blog Categories ----

func (h *Handler) CreateBlogCategory(c *fiber.Ctx) error {
	var req cmsreq.CreateBlogCategory
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	cat, err := h.blogs.CreateCategory(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, cat)
}

func (h *Handler) ListBlogCategories(c *fiber.Ctx) error {
	cats, err := h.blogs.ListCategories(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, cats)
}

func (h *Handler) DeleteBlogCategory(c *fiber.Ctx) error {
	if err := h.blogs.DeleteCategory(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Blogs ----

func (h *Handler) CreateBlog(c *fiber.Ctx) error {
	var req cmsreq.CreateBlog
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	b, err := h.blogs.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, b)
}

func (h *Handler) GetBlog(c *fiber.Ctx) error {
	b, err := h.blogs.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

func (h *Handler) GetBlogBySlug(c *fiber.Ctx) error {
	b, err := h.blogs.GetBySlug(c.Context(), c.Params("website_id"), c.Params("slug"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

func (h *Handler) ListBlogs(c *fiber.Ctx) error {
	filter := cmsreq.ListBlogsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("website_id"); v != "" {
		filter.WebsiteID = &v
	}
	if v := c.Query("blog_category_id"); v != "" {
		filter.BlogCategoryID = &v
	}
	if v := c.Query("status"); v != "" {
		filter.Status = &v
	}
	blogs, err := h.blogs.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, blogs, fiber.Map{"page": filter.Page, "per_page": filter.PerPage})
}

func (h *Handler) UpdateBlog(c *fiber.Ctx) error {
	var req cmsreq.UpdateBlog
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	b, err := h.blogs.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

func (h *Handler) DeleteBlog(c *fiber.Ctx) error {
	if err := h.blogs.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- SEO ----

func (h *Handler) UpsertSEO(c *fiber.Ctx) error {
	var req cmsreq.UpsertSEO
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	meta, err := h.seo.Upsert(c.Context(), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, meta)
}

func (h *Handler) GetSEO(c *fiber.Ctx) error {
	meta, err := h.seo.Get(c.Context(), c.Params("entity_type"), c.Params("entity_id"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, meta)
}
