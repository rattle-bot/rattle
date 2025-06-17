package models

const (
	// Event types
	EventTypeError    = "error"
	EventTypeInfo     = "info"
	EventTypeWarning  = "warning"
	EventTypeSuccess  = "success"
	EventTypeCritical = "critical"

	// Match types
	MatchTypeInclude = "include"
	MatchTypeExclude = "exclude" // Don't send notifications

	// User roles
	RoleAdmin = "admin"
	RoleUser  = "user"

	// Container exclusion types
	ContainerExclusionName  = "name"
	ContainerExclusionImage = "image"
	ContainerExclusionID    = "id"
	ContainerExclusionLabel = "label"
)
