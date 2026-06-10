package handlers

import (
	"kimsha/internal/models"
	"kimsha/internal/repository"
	"kimsha/pkg/password"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct{ repo *repository.UserRepository }

func NewUserHandler(repo *repository.UserRepository) *UserHandler { return &UserHandler{repo: repo} }

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.repo.List(tenantID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, users)
}

type CreateUserInput struct {
	Name     string          `json:"name" validate:"required"`
	NameAm   string          `json:"name_am"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Password string          `json:"password" validate:"required,min=6"`
	Role     models.UserRole `json:"role" validate:"required"`
	PIN      string          `json:"pin"`
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var in CreateUserInput
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	hashed, _ := password.Hash(in.Password)
	var pinHash string
	if in.PIN != "" {
		pinHash, _ = password.Hash(in.PIN)
	}
	u := &models.User{
		ID:       uuid.New(),
		TenantID: tenantID(c),
		Name:     in.Name,
		NameAm:   in.NameAm,
		Email:    in.Email,
		Phone:    in.Phone,
		Password: hashed,
		Role:     in.Role,
		PIN:      pinHash,
		IsActive: true,
	}
	if err := h.repo.Create(u); err != nil {
		return response.Conflict(c, "user already exists")
	}
	return response.Created(c, u)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	u, err := h.repo.Find(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, u)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	u, err := h.repo.Find(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	var in struct {
		Name   string          `json:"name"`
		NameAm string          `json:"name_am"`
		Phone  string          `json:"phone"`
		Role   models.UserRole `json:"role"`
	}
	if err := c.BodyParser(&in); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	u.Name = in.Name
	u.NameAm = in.NameAm
	u.Phone = in.Phone
	if in.Role != "" {
		u.Role = in.Role
	}
	if err := h.repo.Update(u); err != nil {
		return response.Internal(c)
	}
	return response.OK(c, u)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.repo.Delete(id, tenantID(c)); err != nil {
		return response.Internal(c)
	}
	return response.NoContent(c)
}

func (h *UserHandler) SetPIN(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var body struct{ PIN string `json:"pin"` }
	if err := c.BodyParser(&body); err != nil || len(body.PIN) != 4 {
		return response.BadRequest(c, "pin must be 4 digits")
	}
	u, err := h.repo.Find(id, tenantID(c))
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	hashed, _ := password.Hash(body.PIN)
	u.PIN = hashed
	_ = h.repo.Update(u)
	return response.OK(c, fiber.Map{"message": "pin updated"})
}

// ensure *gorm.DB import isn't flagged as unused if compiler complains
var _ = (*gorm.DB)(nil)
