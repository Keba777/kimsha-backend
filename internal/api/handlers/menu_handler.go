package handlers

import (
	"kimsha/internal/api/middleware"
	"kimsha/internal/services"
	"kimsha/pkg/jwt"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MenuHandler struct{ svc *services.MenuService }

func NewMenuHandler(svc *services.MenuService) *MenuHandler { return &MenuHandler{svc: svc} }

func tenantID(c *fiber.Ctx) uuid.UUID {
	claims := c.Locals(middleware.LocalsUserKey).(*jwt.Claims)
	return claims.TenantID
}

// Categories

func (h *MenuHandler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.svc.ListCategories(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, cats)
}

func (h *MenuHandler) CreateCategory(c *fiber.Ctx) error {
	var in services.CategoryInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(in); errs != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}
	cat, err := h.svc.CreateCategory(tenantID(c), in)
	if err != nil {
		return response.Internal(c)
	}
	return response.Created(c, cat)
}

func (h *MenuHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var in services.CategoryInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	cat, err := h.svc.UpdateCategory(id, tenantID(c), in)
	if err != nil {
		return response.NotFound(c, "category not found")
	}
	return response.OK(c, cat)
}

func (h *MenuHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.svc.DeleteCategory(id, tenantID(c)); err != nil {
		return response.Internal(c)
	}
	return response.NoContent(c)
}

// Items

func (h *MenuHandler) ListItems(c *fiber.Ctx) error {
	var catID *uuid.UUID
	if raw := c.Query("category_id"); raw != "" {
		id, _ := uuid.Parse(raw)
		catID = &id
	}
	availableOnly := c.QueryBool("available", false)
	items, err := h.svc.ListItems(tenantID(c), catID, availableOnly)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, items)
}

func (h *MenuHandler) CreateItem(c *fiber.Ctx) error {
	var in services.ItemInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(in); errs != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}
	item, err := h.svc.CreateItem(tenantID(c), in)
	if err != nil {
		return response.Internal(c)
	}
	return response.Created(c, item)
}

func (h *MenuHandler) GetItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	item, err := h.svc.GetItem(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "item not found")
	}
	return response.OK(c, item)
}

func (h *MenuHandler) UpdateItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var in services.ItemInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	item, err := h.svc.UpdateItem(id, tenantID(c), in)
	if err != nil {
		return response.NotFound(c, "item not found")
	}
	return response.OK(c, item)
}

func (h *MenuHandler) DeleteItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.svc.DeleteItem(id, tenantID(c)); err != nil {
		return response.Internal(c)
	}
	return response.NoContent(c)
}

func (h *MenuHandler) ToggleItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	item, err := h.svc.ToggleAvailability(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "item not found")
	}
	return response.OK(c, item)
}
