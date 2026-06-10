package handlers

import (
	"kimsha/internal/api/middleware"
	"kimsha/internal/models"
	"kimsha/internal/services"
	"kimsha/internal/ws"
	"kimsha/pkg/jwt"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
	gows "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type KitchenHandler struct{ svc *services.KitchenService }

func NewKitchenHandler(svc *services.KitchenService) *KitchenHandler {
	return &KitchenHandler{svc: svc}
}

func (h *KitchenHandler) ActiveTickets(c *fiber.Ctx) error {
	tickets, err := h.svc.ActiveTickets(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, tickets)
}

func (h *KitchenHandler) UpdateTicketStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var body struct {
		Status models.TicketStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	claims := c.Locals(middleware.LocalsUserKey).(*jwt.Claims)
	if err := h.svc.UpdateTicketStatus(id, claims.TenantID.String(), body.Status); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"status": body.Status})
}

func (h *KitchenHandler) WebSocket(c *gows.Conn) {
	claims, ok := c.Locals(middleware.LocalsUserKey).(*jwt.Claims)
	if !ok {
		_ = c.Close()
		return
	}
	ws.Default.Register(c, claims.TenantID.String())
}
