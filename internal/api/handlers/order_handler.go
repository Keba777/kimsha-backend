package handlers

import (
	"kimsha/internal/models"
	"kimsha/internal/services"
	"kimsha/pkg/pagination"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct{ svc *services.OrderService }

func NewOrderHandler(svc *services.OrderService) *OrderHandler { return &OrderHandler{svc: svc} }

func (h *OrderHandler) List(c *fiber.Ctx) error {
	p := pagination.Parse(c)
	status := c.Query("status")
	orders, total, err := h.svc.List(tenantID(c), status, p)
	if err != nil {
		return response.Internal(c)
	}
	return response.WithMeta(c, orders, pagination.NewMeta(p, total))
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var in services.CreateOrderInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	order, err := h.svc.Create(tenantID(c), in)
	if err != nil {
		return response.Internal(c)
	}
	return response.Created(c, order)
}

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	order, err := h.svc.Get(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "order not found")
	}
	return response.OK(c, order)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var body struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if err := h.svc.UpdateStatus(id, tenantID(c), body.Status); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"status": body.Status})
}

func (h *OrderHandler) AddItem(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	var in services.AddItemInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(in); errs != nil {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}
	item, err := h.svc.AddItem(orderID, tenantID(c), in)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, item)
}

func (h *OrderHandler) UpdateItemStatus(c *fiber.Ctx) error {
	orderID, _ := uuid.Parse(c.Params("id"))
	itemID, err := uuid.Parse(c.Params("iid"))
	if err != nil {
		return response.BadRequest(c, "invalid item id")
	}
	var body struct {
		Status models.OrderItemStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if err := h.svc.UpdateItem(orderID, itemID, tenantID(c), body.Status); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"status": body.Status})
}

func (h *OrderHandler) RemoveItem(c *fiber.Ctx) error {
	orderID, _ := uuid.Parse(c.Params("id"))
	itemID, err := uuid.Parse(c.Params("iid"))
	if err != nil {
		return response.BadRequest(c, "invalid item id")
	}
	if err := h.svc.RemoveItem(orderID, itemID, tenantID(c)); err != nil {
		return response.Internal(c)
	}
	return response.NoContent(c)
}
