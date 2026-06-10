package response

import "github.com/gofiber/fiber/v2"

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Envelope{Success: true, Data: data})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func WithMeta(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{Success: true, Data: data, Meta: meta})
}

func Err(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(Envelope{Success: false, Error: msg})
}

func BadRequest(c *fiber.Ctx, msg string) error {
	return Err(c, fiber.StatusBadRequest, msg)
}

func Unauthorized(c *fiber.Ctx) error {
	return Err(c, fiber.StatusUnauthorized, "unauthorized")
}

func Forbidden(c *fiber.Ctx) error {
	return Err(c, fiber.StatusForbidden, "forbidden")
}

func NotFound(c *fiber.Ctx, msg string) error {
	return Err(c, fiber.StatusNotFound, msg)
}

func Internal(c *fiber.Ctx) error {
	return Err(c, fiber.StatusInternalServerError, "internal server error")
}

func Conflict(c *fiber.Ctx, msg string) error {
	return Err(c, fiber.StatusConflict, msg)
}
