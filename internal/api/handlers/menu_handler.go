package handlers

import (
	"strconv"

	"kimsha/internal/api/middleware"
	"kimsha/internal/models"
	"kimsha/internal/services"
	"kimsha/pkg/jwt"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MenuHandler struct {
	svc        *services.MenuService
	storageSvc *services.StorageService
}

func NewMenuHandler(svc *services.MenuService, storageSvc *services.StorageService) *MenuHandler {
	return &MenuHandler{svc: svc, storageSvc: storageSvc}
}

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
	in, err := h.parseItemForm(c)
	if err != nil {
		return response.BadRequest(c, err.Error())
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
	in, err := h.parseItemForm(c)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	item, err := h.svc.UpdateItem(id, tenantID(c), in)
	if err != nil {
		return response.NotFound(c, "item not found")
	}
	return response.OK(c, item)
}

// parseItemForm reads multipart form fields and an optional image file.
func (h *MenuHandler) parseItemForm(c *fiber.Ctx) (services.ItemInput, error) {
	var in services.ItemInput

	in.Name = c.FormValue("name")
	in.NameAm = c.FormValue("name_am")
	in.Description = c.FormValue("description")
	in.DescriptionAm = c.FormValue("description_am")
	in.ImageURL = c.FormValue("image_url")

	if raw := c.FormValue("price"); raw != "" {
		p, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return in, fiber.NewError(fiber.StatusBadRequest, "invalid price")
		}
		in.Price = p
	}
	if raw := c.FormValue("prep_time_min"); raw != "" {
		n, _ := strconv.Atoi(raw)
		in.PrepTimeMin = n
	}
	if raw := c.FormValue("sort_order"); raw != "" {
		n, _ := strconv.Atoi(raw)
		in.SortOrder = n
	}
	in.IsFeatured = c.FormValue("is_featured") == "true"
	in.IsAvailable = c.FormValue("is_available") != "false"

	if t := c.FormValue("item_type"); t != "" {
		in.ItemType = models.ItemType(t)
	}
	if raw := c.FormValue("category_id"); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			in.CategoryID = &id
		}
	}

	// Upload image if provided
	if file, err := c.FormFile("image"); err == nil {
		f, err := file.Open()
		if err != nil {
			return in, fiber.NewError(fiber.StatusInternalServerError, "could not read image")
		}
		defer f.Close()
		url, err := h.storageSvc.UploadImage(f, file)
		if err != nil {
			return in, fiber.NewError(fiber.StatusInternalServerError, "image upload failed")
		}
		in.ImageURL = url
	}

	return in, nil
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
