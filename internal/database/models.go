package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID string `gorm:"uniqueIndex"`
	Username   string
	FirstName  string
	Role       string // models.RoleAdmin / RoleUser
}
