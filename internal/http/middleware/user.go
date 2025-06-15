package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/http/handlers"
	"github.com/ilyxenc/rattle/internal/models"
)

func LocatedTelegramId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		tid, ok := claims["id"]
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(handlers.Res{
				Message: "Error check access",
			})
		}

		c.Locals("telegram_id", tid)

		return c.Next()
	}
}

func LocatedUserRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		telegramID := c.Locals("telegram_id")

		db := database.DB

		var userRole string
		if err := db.Model(&models.User{}).Where("telegram_id = ?", telegramID).Pluck("role", &userRole).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(handlers.Res{
				Message: "Not enough rights",
			})
		}

		if userRole == role {
			c.Locals("role", userRole)

			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(handlers.Res{
			Message: "Not enough rights",
		})
	}
}
