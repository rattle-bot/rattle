package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/docker"
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/managers"
	"github.com/ilyxenc/rattle/internal/scanner"
	"github.com/ilyxenc/rattle/internal/telegram"
)

func main() {
	// Load environment configuration from .env or system
	config.Load()
	// Initialize the global logger
	logger.Init("./logs/rattle.log")
	defer logger.Log.Sync() // Flush logs on shutdown

	// Connect to database
	if err := database.Connect("./db/rattle.db"); err != nil {
		logger.Log.Fatal("Failed to connect to database:", err)
	}
	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		logger.Log.Fatal("Failed to migrate database:", err)
	}
	// Initialize data (chat IDs, patterns, exclusions)
	if err := database.Initialize(); err != nil {
		logger.Log.Fatal("Failed to initialize database data:", err)
	}

	// Register observers AFTER DB init but BEFORE logic starts and reload managers
	managers.Init()

	// Initialize Telegram client
	telegram.Init()

	// Log and notify that Rattle has started
	logger.Log.Infof("🚀 Rattle started in %s mode", config.Cfg.Env)
	telegram.Notify(telegram.Notification{
		Type: telegram.NotificationStartedRattle,
	})

	// Create context that cancels on interrupt or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create Docker client and log scan manager
	cli := docker.NewClient()
	manager := scanner.NewLogScanManager(ctx, cli)
	// Start scanning container logs in the background
	go manager.StartAll()

	// Wait until shutdown signal is received
	<-ctx.Done()

	// Trigger shutdown
	manager.StopAll()

	// Log and notify that Rattle is shutting down
	logger.Log.Info("🛑 Shutting down Rattle")
	telegram.Notify(telegram.Notification{
		Type: telegram.NotificationShutDownRattle,
	})
}
