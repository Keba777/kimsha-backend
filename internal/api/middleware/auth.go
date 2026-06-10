package middleware

import (
	"strings"

	"kimsha/pkg/jwt"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
)

const LocalsUserKey = "user"

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return response.Unauthorized(c)
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.Verify(token)
		if err != nil {
			return response.Unauthorized(c)
		}
		c.Locals(LocalsUserKey, claims)
		return c.Next()
	}
}
