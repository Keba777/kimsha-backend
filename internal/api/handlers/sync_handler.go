package handlers

import (
	"kimsha/internal/services"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type SyncHandler struct{ svc *services.SyncService }

func NewSyncHandler(svc *services.SyncService) *SyncHandler { return &SyncHandler{svc: svc} }

func (h *SyncHandler) Push(c *fiber.Ctx) error {
	var req services.PushRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	req.TenantID = tenantID(c)
	results, err := h.svc.Push(tenantID(c), req)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, results)
}

func (h *SyncHandler) MenuSnapshot(c *fiber.Ctx) error {
	snapshot, err := h.svc.MenuSnapshot(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, snapshot)
}

func (h *SyncHandler) TableSnapshot(c *fiber.Ctx) error {
	tables, err := h.svc.TableSnapshot(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, tables)
}
