package managers

import (
	"sync"
	"time"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/logger"
	"gorm.io/gorm"
)

// ChangeWatcher tracks a table and its timestamp fields to detect changes
type ChangeWatcher struct {
	Table       string    // Table name to monitor
	Fields      []string  // Fields like "updated_at", "deleted_at"
	ReloadFunc  func()    // Function to call when a change is detected
	LastChecked time.Time // Timestamp of the last check
}

// ChangeWatcherManager holds and manages all registered watchers
type ChangeWatcherManager struct {
	watchers []*ChangeWatcher // All watchers
	ticker   *time.Ticker     // Shared polling ticker
	stopChan chan struct{}    // Channel to stop polling
	mu       sync.Mutex       // Mutex for thread-safe access
}

var watcherManager = &ChangeWatcherManager{
	watchers: make([]*ChangeWatcher, 0),
	stopChan: make(chan struct{}),
}

// AddWatcher registers a new table + field(s) watcher and reload function
func AddWatcher(table string, fields []string, reload func()) {
	watcher := &ChangeWatcher{
		Table:       table,
		Fields:      fields,
		ReloadFunc:  reload,
		LastChecked: time.Now().Add(-5 * time.Minute), // Start a bit in the past
	}
	watcherManager.mu.Lock()
	watcherManager.watchers = append(watcherManager.watchers, watcher)
	watcherManager.mu.Unlock()
}

// StartWatchers begins periodic polling of all watchers
func StartWatchers(interval time.Duration) {
	watcherManager.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-watcherManager.stopChan:
				return
			case <-watcherManager.ticker.C:
				watcherManager.runChecks()
			}
		}
	}()
}

// StopWatchers gracefully stops the polling goroutine
func StopWatchers() {
	close(watcherManager.stopChan)
}

// runChecks executes checks for all registered watchers
func (m *ChangeWatcherManager) runChecks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, w := range m.watchers {
		for _, field := range w.Fields {
			db := database.DB

			var count int64

			// Query rows with updated_at > last check
			err := db.Session(&gorm.Session{}).
				Table(w.Table).
				Where(field+" > ?", w.LastChecked).
				Count(&count).Error

			if err != nil {
				logger.Log.Warnf("Watcher error on table %s: %v", w.Table, err)
				continue
			}

			if count > 0 {
				logger.Log.Infof("Detected change in %s (%s), reloading...", w.Table, field)
				w.LastChecked = time.Now()
				w.ReloadFunc()
				break // One match is enough to trigger reload
			}
		}
	}
}
