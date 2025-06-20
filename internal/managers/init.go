package managers

import (
	"time"

	"github.com/ilyxenc/rattle/internal/logger"
)

func Init() {
	// Initial reload of data for managers
	if err := Chats.Reload(); err != nil {
		logger.Log.Fatalf("Failed to load chat IDs: %v", err)
	}
	if err := Containers.Reload(); err != nil {
		logger.Log.Fatalf("Failed to load container exclusions: %v", err)
	}
	if err := Logs.Reload(); err != nil {
		logger.Log.Fatalf("Failed to load logs exclusions: %v", err)
	}

	// Register table watchers (no duplicate interval)
	AddWatcher("log_exclusions", []string{"updated_at", "deleted_at"}, func() {
		if err := Logs.Reload(); err != nil {
			logger.Log.Warnf("Failed to reload log exclusions: %v", err)
		}
	})
	AddWatcher("container_exclusions", []string{"updated_at", "deleted_at"}, func() {
		if err := Containers.Reload(); err != nil {
			logger.Log.Warnf("Failed to reload container exclusions: %v", err)
		}
	})
	AddWatcher("chats", []string{"updated_at", "deleted_at"}, func() {
		if err := Chats.Reload(); err != nil {
			logger.Log.Warnf("Failed to reload chat IDs: %v", err)
		}
	})
	AddWatcher("mode", []string{"updated_at", "deleted_at"}, func() {
		if err := Mode.Reload(); err != nil {
			logger.Log.Warnf("Failed to reload mode: %v", err)
		}
	})

	// Start polling every 15 seconds
	StartWatchers(15 * time.Second)
}
