package models

import (
	"gorm.io/gorm"
)

type ContainerExclusion struct {
	gorm.Model
	Type  string `json:"type"`  // models.ContainerExclusionName / Image / ID
	Value string `json:"value"` // Include value to exclude from all
}
