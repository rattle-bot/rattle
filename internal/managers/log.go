package managers

import (
	"regexp"
	"strings"
	"sync"

	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/logger"
	"github.com/ilyxenc/rattle/internal/models"
)

type LogManager struct {
	mu      sync.RWMutex
	cache   map[string][]*regexp.Regexp // map[EventType] = compiled regex patterns
	exclude []*regexp.Regexp            // exclude patterns (no event type)
}

// Logs is the global log manager instance
var Logs = &LogManager{
	cache: make(map[string][]*regexp.Regexp),
}

// Reload fetches log patterns from DB and compiles them
func (lm *LogManager) Reload() error {
	var patterns []models.LogExclusion

	if err := database.DB.Find(&patterns).Error; err != nil {
		return err
	}

	newCache := make(map[string][]*regexp.Regexp)
	newExclude := make([]*regexp.Regexp, 0, len(patterns))

	for _, p := range patterns {
		pattern := strings.TrimSpace(p.Pattern)
		if pattern == "" {
			continue
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			logger.Log.Warnf("Invalid regex pattern: %s", pattern)
			continue
		}

		if p.MatchType == models.MatchTypeExclude {
			newExclude = append(newExclude, regex)
		} else {
			eventType := strings.ToLower(p.EventType)
			newCache[eventType] = append(newCache[eventType], regex)
		}
	}

	lm.mu.Lock()
	lm.cache = newCache
	lm.exclude = newExclude
	lm.mu.Unlock()

	return nil
}

// Include returns compiled patterns for the given event type
func (lm *LogManager) Include(eventType string) []*regexp.Regexp {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.cache[strings.ToLower(eventType)]
}

// Exclude returns compiled exclusion patterns
func (lm *LogManager) Exclude() []*regexp.Regexp {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.exclude
}

// KnownEventTypes returns a list of all event types that exist in memory
func (lm *LogManager) KnownEventTypes() []string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	eventTypes := make([]string, 0, len(lm.cache))
	for et := range lm.cache {
		eventTypes = append(eventTypes, et)
	}
	return eventTypes
}

// OnLogChanged is triggered by the model hook
func (lm *LogManager) OnLogChanged() {
	if err := lm.Reload(); err != nil {
		logger.Log.Errorf("Failed to reload log patterns: %v", err)
	}
}
