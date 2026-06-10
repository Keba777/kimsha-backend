package middleware

import (
	"kimsha/pkg/jwt"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(LocalsUserKey).(*jwt.Claims)
		if !ok {
			return response.Unauthorized(c)
		}
		if !allowed[claims.Role] {
			return response.Forbidden(c)
		}
		return c.Next()
	}
}
