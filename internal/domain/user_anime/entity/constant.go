package entity

// Status is user anime status.
type Status string

// Available user anime status.
const (
	StatusWatching  Status = "WATCHING"
	StatusCompleted Status = "COMPLETED"
	StatusOnHold    Status = "ON_HOLD"
	StatusDropped   Status = "DROPPED"
	StatusPlanned   Status = "PLANNED"
)

// Priority is user anime priority.
type Priority string

// Available user anime priority.
const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

// RewatchValue is user anime rewatch value.
type RewatchValue string

// Available user anime rewatch value.
const (
	RewatchValueVeryLow  RewatchValue = "VERY_LOW"
	RewatchValueLow      RewatchValue = "LOW"
	RewatchValueMedium   RewatchValue = "MEDIUM"
	RewatchValueHigh     RewatchValue = "HIGH"
	RewatchValueVeryHigh RewatchValue = "VERY_HIGH"
)
