package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/http"
	"github.com/ilyxenc/rattle/internal/http/router"
	"github.com/ilyxenc/rattle/internal/logger"
)

func main() {
	// Load environment configuration from .env or system
	config.Load()

	// Initialize the global logger
	logger.Init("./logs/rattle.log")
	defer logger.Log.Sync() // Flush logs on shutdown

	http.InitKeys()

	// Connect to database
	if err := database.Connect(); err != nil {
		logger.Log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB body size
	})

	app.Use(cors.New())

	router.SetupRoutes(app)

	// Port from .env
	logger.Log.Fatal(app.Listen(fmt.Sprintf(":%v", config.Cfg.Fiber.Port)))
}
