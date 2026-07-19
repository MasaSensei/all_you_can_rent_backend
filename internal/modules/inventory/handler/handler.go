package handler

import (
	"github.com/gofiber/fiber/v2"

	invreq "rentos/internal/modules/inventory/dto/request"
	"rentos/internal/modules/inventory/entity"
	"rentos/internal/modules/inventory/service"
	apiresponse "rentos/pkg/response"
	"rentos/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups the inventory module's HTTP handlers.
type Handler struct {
	categories service.CategoryService
	templates  service.AssetTemplateService
	assets     service.AssetService
	validate   *validator.Validate
}

func New(
	categories service.CategoryService,
	templates service.AssetTemplateService,
	assets service.AssetService,
	v *validator.Validate,
) *Handler {
	return &Handler{categories: categories, templates: templates, assets: assets, validate: v}
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

// ---- Categories ----

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var req invreq.CreateCategory
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	cat, err := h.categories.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, cat)
}

func (h *Handler) GetCategory(c *fiber.Ctx) error {
	cat, err := h.categories.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, cat)
}

func (h *Handler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.categories.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, buildCategoryTree(cats))
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	var req invreq.UpdateCategory
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	cat, err := h.categories.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, cat)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.categories.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Asset Templates ----

func (h *Handler) CreateAssetTemplate(c *fiber.Ctx) error {
	var req invreq.CreateAssetTemplate
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	t, err := h.templates.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, t)
}

func (h *Handler) GetAssetTemplate(c *fiber.Ctx) error {
	t, err := h.templates.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, t)
}

func (h *Handler) ListAssetTemplates(c *fiber.Ctx) error {
	ts, err := h.templates.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, ts)
}

func (h *Handler) UpdateAssetTemplate(c *fiber.Ctx) error {
	var req invreq.UpdateAssetTemplate
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	t, err := h.templates.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, t)
}

func (h *Handler) DeleteAssetTemplate(c *fiber.Ctx) error {
	if err := h.templates.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Assets ----

func (h *Handler) CreateAsset(c *fiber.Ctx) error {
	var req invreq.CreateAsset
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	a, err := h.assets.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, a)
}

func (h *Handler) GetAsset(c *fiber.Ctx) error {
	a, err := h.assets.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, a)
}

func (h *Handler) ListAssets(c *fiber.Ctx) error {
	filter := invreq.ListAssetsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("category_id"); v != "" {
		filter.CategoryID = &v
	}
	if v := c.Query("status"); v != "" {
		filter.Status = &v
	}
	if v := c.Query("search"); v != "" {
		filter.Search = &v
	}

	assets, err := h.assets.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, assets, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) UpdateAsset(c *fiber.Ctx) error {
	var req invreq.UpdateAsset
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	a, err := h.assets.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, a)
}

func (h *Handler) DeleteAsset(c *fiber.Ctx) error {
	if err := h.assets.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Asset Images ----

func (h *Handler) AddAssetImage(c *fiber.Ctx) error {
	var body struct {
		URL       string  `json:"url" validate:"required,url"`
		AltText   *string `json:"alt_text"`
		IsPrimary bool    `json:"is_primary"`
		SortOrder int     `json:"sort_order"`
	}
	if err := c.BodyParser(&body); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(body); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	img := &entity.AssetImage{
		URL: body.URL, AltText: body.AltText,
		IsPrimary: body.IsPrimary, SortOrder: body.SortOrder,
	}
	result, err := h.assets.AddImage(c.Context(), tenantID(c), c.Params("id"), userID(c), img)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, result)
}

func (h *Handler) DeleteAssetImage(c *fiber.Ctx) error {
	if err := h.assets.DeleteImage(c.Context(), tenantID(c), c.Params("id"), c.Params("image_id"), userID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Asset Documents ----

func (h *Handler) AddAssetDocument(c *fiber.Ctx) error {
	var body struct {
		Title    string  `json:"title" validate:"required,max=255"`
		FileURL  string  `json:"file_url" validate:"required,url"`
		FileType *string `json:"file_type"`
		FileSize *int    `json:"file_size"`
	}
	if err := c.BodyParser(&body); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(body); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	doc := &entity.AssetDocument{
		Title: body.Title, FileURL: body.FileURL,
		FileType: body.FileType, FileSize: body.FileSize,
	}
	result, err := h.assets.AddDocument(c.Context(), tenantID(c), c.Params("id"), userID(c), doc)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, result)
}

func (h *Handler) DeleteAssetDocument(c *fiber.Ctx) error {
	if err := h.assets.DeleteDocument(c.Context(), tenantID(c), c.Params("id"), c.Params("doc_id"), userID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Availability ----

func (h *Handler) BlockAvailability(c *fiber.Ctx) error {
	var req invreq.BlockAvailability
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	av, err := h.assets.BlockAvailability(c.Context(), tenantID(c), c.Params("id"), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, av)
}

func (h *Handler) ListAvailability(c *fiber.Ctx) error {
	list, err := h.assets.ListAvailability(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, list)
}

// ---- helpers ----

// buildCategoryTree converts a flat list into a parent→children hierarchy.
// Categories without a ParentID become top-level nodes.
func buildCategoryTree(cats []entity.Category) []entity.Category {
	index := make(map[string]*entity.Category, len(cats))
	for i := range cats {
		index[cats[i].ID] = &cats[i]
	}
	var roots []entity.Category
	for i := range cats {
		c := &cats[i]
		if c.ParentID == nil {
			roots = append(roots, *c)
		}
	}
	return roots
}
