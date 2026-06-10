package handlers

import (
	"kimsha/internal/models"
	"kimsha/internal/services"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TableHandler struct{ svc *services.TableService }

func NewTableHandler(svc *services.TableService) *TableHandler { return &TableHandler{svc: svc} }

func (h *TableHandler) List(c *fiber.Ctx) error {
	tables, err := h.svc.List(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, tables)
}

func (h *TableHandler) Create(c *fiber.Ctx) error {
	var in services.TableInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(in); errs != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}
	t, err := h.svc.Create(tenantID(c), in)
	if err != nil {
		return response.Internal(c)
	}
	return response.Created(c, t)
}

func (h *TableHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var in services.TableInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	t, err := h.svc.Update(id, tenantID(c), in)
	if err != nil {
		return response.NotFound(c, "table not found")
	}
	return response.OK(c, t)
}

func (h *TableHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.svc.Delete(id, tenantID(c)); err != nil {
		return response.Internal(c)
	}
	return response.NoContent(c)
}

func (h *TableHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var body struct {
		Status models.TableStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if err := h.svc.UpdateStatus(id, tenantID(c), body.Status); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"status": body.Status})
}
