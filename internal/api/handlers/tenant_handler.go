package handlers

import (
	"kimsha/internal/api/middleware"
	"kimsha/internal/models"
	"kimsha/pkg/jwt"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TenantHandler struct{ db *gorm.DB }

func NewTenantHandler(db *gorm.DB) *TenantHandler { return &TenantHandler{db: db} }

func (h *TenantHandler) Get(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*jwt.Claims)
	var tenant models.Tenant
	if err := h.db.First(&tenant, "id = ?", claims.TenantID).Error; err != nil {
		return response.NotFound(c, "tenant not found")
	}
	return response.OK(c, tenant)
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*jwt.Claims)
	var tenant models.Tenant
	if err := h.db.First(&tenant, "id = ?", claims.TenantID).Error; err != nil {
		return response.NotFound(c, "tenant not found")
	}
	var in struct {
		Name          string  `json:"name"`
		NameAm        string  `json:"name_am"`
		Phone         string  `json:"phone"`
		Address       string  `json:"address"`
		TaxRate       float64 `json:"tax_rate"`
		ServiceCharge float64 `json:"service_charge"`
	}
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if in.Name != "" {
		tenant.Name = in.Name
	}
	tenant.NameAm = in.NameAm
	tenant.Phone = in.Phone
	tenant.Address = in.Address
	tenant.TaxRate = in.TaxRate
	tenant.ServiceCharge = in.ServiceCharge
	h.db.Save(&tenant)
	return response.OK(c, tenant)
}
