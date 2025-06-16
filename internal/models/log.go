package models

import (
	"gorm.io/gorm"
)

type LogExclusion struct {
	gorm.Model
	Pattern   string `json:"pattern"`    // regex-pattern
	MatchType string `json:"match_type"` // models.MatchTypeInclude / MatchTypeExclude
	EventType string `json:"event_type"` // models.EventTypeError / etc
}
