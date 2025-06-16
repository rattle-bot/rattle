package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID string `gorm:"uniqueIndex" json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	Role       string `json:"role"` // models.RoleAdmin / RoleUser
	Active     bool   `json:"active"`
}
