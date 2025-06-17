package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Heartbeat(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "ok",
	})
}
