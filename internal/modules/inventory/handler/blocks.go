package handler

import (
	"github.com/gofiber/fiber/v2"

	invreq "rentos-backend/internal/modules/inventory/dto/request"
	apiresponse "rentos-backend/pkg/response"
)

// GET /assets/:id/blocks
func (h *Handler) ListBlocks(c *fiber.Ctx) error {
	list, err := h.assets.ListAvailability(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, list)
}

// POST /assets/:id/blocks
func (h *Handler) CreateBlock(c *fiber.Ctx) error {
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
