package managers

import (
	"sync"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

// ModeManager handles the current container filtering mode in memory
type ModeManager struct {
	mu    sync.RWMutex
	value string
}

// Mode is the global instance
var Mode = &ModeManager{}

// Reload fetches the filtering mode from the database
func (m *ModeManager) Reload() error {
	var mode models.Mode

	if err := database.DB.First(&mode).Error; err != nil {
		return err
	}

	m.mu.Lock()
	m.value = mode.Value
	m.mu.Unlock()

	return nil
}

// Get returns the currently cached filtering mode
func (m *ModeManager) Get() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.value
}
