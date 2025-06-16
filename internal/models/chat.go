package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ChatID string `gorm:"uniqueIndex" json:"chat_id"`
	Send   bool   `gorm:"default:true" json:"send"` // Send notifications only if true
}

// Callback interface
type ChatObserver interface {
	OnChatChanged()
}

var chatObserver ChatObserver // Registered externally

func RegisterChatObserver(o ChatObserver) {
	chatObserver = o
}

func (c *Chat) AfterCreate(tx *gorm.DB) error {
	if chatObserver != nil {
		chatObserver.OnChatChanged()
	}
	return nil
}

func (c *Chat) AfterUpdate(tx *gorm.DB) error {
	if chatObserver != nil {
		chatObserver.OnChatChanged()
	}
	return nil
}

func (c *Chat) AfterDelete(tx *gorm.DB) error {
	if chatObserver != nil {
		chatObserver.OnChatChanged()
	}
	return nil
}
