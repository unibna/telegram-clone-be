package middleware

import (
	"chat-app/internal/constants"
	"chat-app/pkg/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(secret string) fiber.Handler {
	log.Println("-----> Check auth")

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			auth = c.Query("access_token")
			if auth == "" {
				return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
					"code":    constants.StatusUnauthorized,
					"message": constants.ErrUnauthorized,
				})
			}
		}

		// Kiểm tra và tách "Bearer "
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"code":    constants.StatusUnauthorized,
				"message": constants.ErrTokenInvalid,
			})
		}

		token := parts[1]
		userID, err := utils.ValidateToken(token, secret)
		if err != nil {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"code":    constants.StatusUnauthorized,
				"message": constants.ErrTokenInvalid,
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
