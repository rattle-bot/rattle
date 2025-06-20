package managers

import (
	"strings"
	"sync"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
	"golang.org/x/exp/slices"
)

// ContainerManager handles exclusion logic and in-memory cache
type ContainerManager struct {
	mu    sync.RWMutex                   // Mutex to protect concurrent access to the cache
	cache map[string]map[string][]string // // cache[mode][type] = []values. Key = Type ("name", "image", "id", "label"), Value = list of excluded strings
}

// Containers is a globally accessible instance of ContainerManager
var Containers = &ContainerManager{
	cache: make(map[string]map[string][]string),
}

// Reload fetches excluded containers from the database and stores them in memory
func (m *ContainerManager) Reload() error {
	var all []models.Container

	if err := database.DB.Find(&all).Error; err != nil {
		return err
	}

	newCache := map[string]map[string][]string{
		models.Blacklist: {},
		models.Whitelist: {},
	}

	for _, c := range all {
		mode := strings.ToLower(c.Mode)
		t := strings.ToLower(c.Type)

		if _, ok := newCache[mode]; !ok {
			newCache[mode] = map[string][]string{}
		}
		newCache[mode][t] = append(newCache[mode][t], strings.ToLower(c.Value))
	}

	m.mu.Lock()
	m.cache = newCache
	m.mu.Unlock()

	return nil
}

// All returns a copy of containers with selected `type` (`t`) stored in memory
func (m *ContainerManager) All(t, mode string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	byMode, ok := m.cache[mode]
	if !ok {
		return nil
	}

	return slices.Clone(byMode[strings.ToLower(t)])
}
