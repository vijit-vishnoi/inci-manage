package models

import "time"

type WorkItem struct {
	ID          int       `json:"id"`
	ComponentID string    `json:"component_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    int       `json:"severity"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
