package handlers

import (
	"kimsha/internal/services"
	"kimsha/pkg/response"
	"kimsha/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc           *services.AuthService
	allowSelfReg  bool
}

func NewAuthHandler(svc *services.AuthService, allowSelfReg bool) *AuthHandler {
	return &AuthHandler{svc: svc, allowSelfReg: allowSelfReg}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	if !h.allowSelfReg {
		return response.Err(c, fiber.StatusForbidden, "self-registration is disabled")
	}
	var input services.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(input); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}
	res, err := h.svc.Register(input)
	if err != nil {
		return response.Conflict(c, err.Error())
	}
	return response.Created(c, res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input services.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	if errs := validator.Validate(input); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}
	res, err := h.svc.Login(input)
	if err != nil {
		return response.Err(c, fiber.StatusUnauthorized, err.Error())
	}
	return response.OK(c, res)
}

func (h *AuthHandler) PINLogin(c *fiber.Ctx) error {
	var input services.PINLoginInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "invalid body")
	}
	res, err := h.svc.PINLogin(input)
	if err != nil {
		return response.Err(c, fiber.StatusUnauthorized, err.Error())
	}
	return response.OK(c, res)
}
