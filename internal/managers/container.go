package managers

import (
	"strings"
	"sync"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/models"
	"golang.org/x/exp/slices"
)

// ContainerManager handles exclusion logic and in-memory cache
type ContainerManager struct {
	mu    sync.RWMutex        // Mutex to protect concurrent access to the cache
	cache map[string][]string // Key = Type ("name", "image", "id"), Value = list of excluded strings
}

// Containers is a globally accessible instance of ContainerManager
var Containers = &ContainerManager{
	cache: make(map[string][]string),
}

// Reload fetches excluded containers from the database and stores them in memory
func (m *ContainerManager) Reload() error {
	var all []models.ContainerExclusion

	if err := database.DB.Find(&all).Error; err != nil {
		return err
	}

	newCache := map[string][]string{}
	for _, e := range all {
		t := strings.ToLower(e.Type)
		newCache[t] = append(newCache[t], strings.ToLower(e.Value))
	}

	m.mu.Lock()
	m.cache = newCache
	m.mu.Unlock()

	return nil
}

// All returns a copy of containers with selected `type` (`t`) stored in memory
func (m *ContainerManager) All(t string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return slices.Clone(m.cache[strings.ToLower(t)])
}

// OnContainerChanged is called automatically by GORM hooks when a ContainerExclusion is created, updated, or deleted. It refreshes the in-memory cache
func (cm *ContainerManager) OnContainerChanged() {
	if err := cm.Reload(); err != nil {
		logger.Log.Errorf("Failed to reload containers after change: %v", err)
	}
}
