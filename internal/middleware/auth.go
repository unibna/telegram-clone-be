package middleware

import (
	"chat-app/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"chat-app/internal/constants"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"code": constants.StatusUnauthorized,
				"message": constants.ErrUnauthorized,
			})
		}

		userID, err := utils.ValidateToken(token, secret)
		if err != nil {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"code": constants.StatusUnauthorized,
				"message": constants.ErrTokenInvalid,
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
} 