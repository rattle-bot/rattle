package models

import (
	"gorm.io/gorm"
)

type Mode struct {
	gorm.Model
	Value string `gorm:"default:blacklist" json:"value"` // Current filtering mode models.Blacklist / Whitelist
}
