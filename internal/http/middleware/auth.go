package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyxenc/rattle/internal/http"
	"github.com/ilyxenc/rattle/internal/http/handlers"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return http.PublicKey, nil
		},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		if err.Error() == "missing or malformed JWT" {
			return c.Status(fiber.StatusBadRequest).JSON(handlers.Res{
				Message: "Missing token",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(handlers.Res{
			Message: "Invalid token",
		})
	}

	return c.Next()
}
