package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ilyxenc/rattle/internal/http/handlers"
	mw "github.com/ilyxenc/rattle/internal/http/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	auth := api.Group("/auth")
	auth.Post("/", handlers.AuthTelegram)

	user := api.Group("/user")
	user.Post("/new", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.CreateUser)

}
