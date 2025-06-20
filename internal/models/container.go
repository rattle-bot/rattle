package models

import (
	"gorm.io/gorm"
)

type Container struct {
	gorm.Model
	Type  string `json:"type"`                          // models.ContainerName / Image / ID
	Value string `json:"value"`                         // Include value to exclude from all
	Mode  string `gorm:"default:blacklist" json:"mode"` // models.Blacklist / Whitelist
}
