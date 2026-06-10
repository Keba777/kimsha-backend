package handlers

import (
	"fmt"
	"strings"
	"time"

	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/pkg/jwt"
	"kimsha/pkg/password"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	repo   *repository.AdminRepository
	expiry time.Duration
}

func NewAdminHandler(repo *repository.AdminRepository, expiry time.Duration) *AdminHandler {
	return &AdminHandler{repo: repo, expiry: expiry}
}

func (h *AdminHandler) Login(c *fiber.Ctx) error {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	admin, err := h.repo.FindSuperAdmin(in.Email)
	if err != nil {
		return response.Err(c, fiber.StatusUnauthorized, "invalid credentials")
	}
	if !password.Compare(admin.Password, in.Password) {
		return response.Err(c, fiber.StatusUnauthorized, "invalid credentials")
	}
	token, err := jwt.Sign(admin.ID, uuid.Nil, string(admin.Role), admin.Name, h.expiry)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"token": token, "user": admin})
}

func (h *AdminHandler) Stats(c *fiber.Ctx) error {
	stats, err := h.repo.GetStats()
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, stats)
}

func (h *AdminHandler) ListTenants(c *fiber.Ctx) error {
	tenants, err := h.repo.ListTenants()
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, tenants)
}

type CreateTenantInput struct {
	TenantName string `json:"tenant_name"`
	TenantSlug string `json:"tenant_slug"`
	OwnerName  string `json:"owner_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Plan       string `json:"plan"`
}

func (h *AdminHandler) CreateTenant(c *fiber.Ctx) error {
	var in CreateTenantInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if in.TenantName == "" || in.TenantSlug == "" || in.OwnerName == "" || in.Email == "" || in.Password == "" {
		return response.BadRequest(c, "tenant_name, tenant_slug, owner_name, email and password are required")
	}
	plan := in.Plan
	if plan == "" {
		plan = "free"
	}
	tenant := &models.Tenant{
		ID:   uuid.New(),
		Name: in.TenantName,
		Slug: strings.ToLower(in.TenantSlug),
		Plan: plan,
	}
	hashed, err := password.Hash(in.Password)
	if err != nil {
		return response.Internal(c)
	}
	owner := &models.User{
		ID:       uuid.New(),
		Name:     in.OwnerName,
		Email:    in.Email,
		Password: hashed,
		Role:     models.RoleOwner,
		IsActive: true,
	}
	if err := h.repo.CreateTenantWithOwner(tenant, owner); err != nil {
		return response.Conflict(c, fmt.Sprintf("could not create tenant: %v", err))
	}
	return response.Created(c, fiber.Map{"tenant": tenant, "owner": owner})
}

func (h *AdminHandler) GetTenant(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	detail, err := h.repo.GetTenantDetail(id)
	if err != nil {
		return response.NotFound(c, "tenant not found")
	}
	return response.OK(c, detail)
}

func (h *AdminHandler) SetTenantStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var in struct {
		Active bool `json:"active"`
	}
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if err := h.repo.SetTenantStatus(id, in.Active); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"active": in.Active})
}
