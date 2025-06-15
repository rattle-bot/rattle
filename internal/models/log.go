package models

import "gorm.io/gorm"

type LogExclusion struct {
	gorm.Model
	Pattern   string // regex-pattern
	MatchType string // models.MatchTypeInclude / MatchTypeExclude
	EventType string // models.EventTypeError / etc
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
