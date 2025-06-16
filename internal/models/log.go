package models

import "gorm.io/gorm"

type LogExclusion struct {
	gorm.Model
	Pattern   string `json:"pattern"` // regex-pattern
	MatchType string `json:"match_type"` // models.MatchTypeInclude / MatchTypeExclude
	EventType string `json:"event_type"` // models.EventTypeError / etc
}

// Callback interface
type LogObserver interface {
	OnLogChanged()
}

var logObserver LogObserver // Registered externally

func RegisterLogObserver(o LogObserver) {
	logObserver = o
}

func (c *LogExclusion) AfterCreate(tx *gorm.DB) error {
	if logObserver != nil {
		logObserver.OnLogChanged()
	}
	return nil
}

func (c *LogExclusion) AfterUpdate(tx *gorm.DB) error {
	if logObserver != nil {
		logObserver.OnLogChanged()
	}
	return nil
}

func (c *LogExclusion) AfterDelete(tx *gorm.DB) error {
	if logObserver != nil {
		logObserver.OnLogChanged()
	}
	return nil
}
