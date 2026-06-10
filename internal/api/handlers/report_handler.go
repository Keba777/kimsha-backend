package handlers

import (
	"time"

	"kimsha/internal/services"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct{ svc *services.ReportService }

func NewReportHandler(svc *services.ReportService) *ReportHandler { return &ReportHandler{svc: svc} }

func (h *ReportHandler) Daily(c *fiber.Ctx) error {
	date := time.Now()
	if raw := c.Query("date"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			date = t
		}
	}
	data, err := h.svc.Daily(tenantID(c), date)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, data)
}

func (h *ReportHandler) TopItems(c *fiber.Ctx) error {
	from, to := dateRange(c)
	items, err := h.svc.TopItems(tenantID(c), from, to)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, items)
}

func (h *ReportHandler) HourlySales(c *fiber.Ctx) error {
	date := time.Now()
	if raw := c.Query("date"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			date = t
		}
	}
	sales, err := h.svc.HourlySales(tenantID(c), date)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, sales)
}

func (h *ReportHandler) WaiterStats(c *fiber.Ctx) error {
	from, to := dateRange(c)
	stats, err := h.svc.WaiterStats(tenantID(c), from, to)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, stats)
}

func dateRange(c *fiber.Ctx) (time.Time, time.Time) {
	from := time.Now().Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)
	if raw := c.Query("from"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			from = t
		}
	}
	if raw := c.Query("to"); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			to = t.Add(24 * time.Hour)
		}
	}
	return from, to
}
