package managers

import (
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/models"
)

// InitObservers registers all in-memory managers as observers for their respective models.
// This enables automatic reloading of cached data when the corresponding database tables change.
// For example, when a Chat, ContainerExclusion or LogExclusion is created/updated/deleted, the corresponding manager will automatically refresh its in-memory cache
func initObservers() {
	models.RegisterChatObserver(Chats)           // Reload chat IDs on changes
	models.RegisterContainerObserver(Containers) // Reload container exclusions on changes
	models.RegisterLogObserver(Logs)             // Reload logs exclusions on changes
}

func Init() {
	// Register observers
	initObservers()

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
}
