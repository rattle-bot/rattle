package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ilyxenc/rattle/internal/http/handlers"
	mw "github.com/ilyxenc/rattle/internal/http/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/heartbeat", handlers.Heartbeat)

	auth := api.Group("/auth")
	auth.Post("/", handlers.AuthTelegram)

	user := api.Group("/user")
	user.Post("/new", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.CreateUser)
	user.Get("/", mw.Protected(), mw.LocatedTelegramId(), handlers.GetMe)
	user.Get("/list", mw.Protected(), handlers.ListUsers)
	user.Delete("/:telegram_id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.DeleteUser)

	chat := api.Group("/chat")
	chat.Post("/new", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.CreateChat)
	chat.Get("/list", mw.Protected(), handlers.ListChats)
	chat.Patch("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.UpdateChat)
	chat.Delete("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.DeleteChat)

	container := api.Group("/container")
	container.Post("/new", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.CreateContainer)
	container.Get("/list", mw.Protected(), handlers.ListContainers)
	container.Patch("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.UpdateContainer)
	container.Delete("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.DeleteContainer)

	log := api.Group("/log")
	log.Post("/new", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.CreateLog)
	log.Get("/list", mw.Protected(), handlers.ListLog)
	log.Patch("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.UpdateLog)
	log.Delete("/:id", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.DeleteLog)

	mode := api.Group("/mode")
	mode.Patch("/", mw.Protected(), mw.LocatedTelegramId(), mw.LocatedUserRole("admin"), handlers.UpdateFilteringMode)
}
