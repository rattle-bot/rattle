package models

import (
	"gorm.io/gorm"
)

type ContainerExclusion struct {
	gorm.Model
	Type  string `json:"type"` // models.ContainerExclusionName / Image / ID
	Value string `json:"value"` // Include value to exclude from all
}

// Callback interface
type ContainerObserver interface {
	OnContainerChanged()
}

var containerObserver ContainerObserver

func RegisterContainerObserver(o ContainerObserver) {
	containerObserver = o
}

func (ce *ContainerExclusion) AfterCreate(tx *gorm.DB) error {
	if containerObserver != nil {
		containerObserver.OnContainerChanged()
	}
	return nil
}

func (ce *ContainerExclusion) AfterUpdate(tx *gorm.DB) error {
	if containerObserver != nil {
		containerObserver.OnContainerChanged()
	}
	return nil
}

func (ce *ContainerExclusion) AfterDelete(tx *gorm.DB) error {
	if containerObserver != nil {
		containerObserver.OnContainerChanged()
	}
	return nil
}
