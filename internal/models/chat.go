package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ChatID string `gorm:"uniqueIndex" json:"chat_id"`
	Send   bool   `gorm:"default:true" json:"send"` // Send notifications only if true
}
