package models

import "time"

// Signal represents the incoming high-volume JSON error payload.
type Signal struct {
	ComponentID string                 `json:"component_id" binding:"required"`
	ErrorCode   string                 `json:"error_code" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}
