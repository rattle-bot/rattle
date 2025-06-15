package managers

import (
	"sync"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/models"
	"golang.org/x/exp/slices"
)

// ChatManager is responsible for managing active chat IDs in memory and keeping them in sync with the database
type ChatManager struct {
	mu      sync.RWMutex // Read-write mutex to protect concurrent access
	chatIDs []string     // Cached list of active chat IDs
}

// Chats is a globally accessible instance of ChatManager
var Chats = &ChatManager{}

// Reload fetches active chat IDs from the database and stores them in memory
func (cm *ChatManager) Reload() error {
	var chats []models.Chat

	if err := database.DB.Where("send = ?", true).Find(&chats).Error; err != nil {
		return err
	}

	var ids []string
	for _, chat := range chats {
		ids = append(ids, chat.ChatID)
	}

	cm.mu.Lock()
	cm.chatIDs = ids
	cm.mu.Unlock()

	return nil
}

// All returns a copy of all active chat IDs stored in memory
func (cm *ChatManager) All() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return slices.Clone(cm.chatIDs)
}

// OnChatChanged is called automatically by GORM hooks when a Chat is created, updated, or deleted. It refreshes the in-memory cache
func (cm *ChatManager) OnChatChanged() {
	if err := cm.Reload(); err != nil {
		logger.Log.Errorf("Failed to reload chat IDs after change: %v", err)
	}
}
