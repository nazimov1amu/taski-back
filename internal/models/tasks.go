package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	PlannedAt   string    `json:"planned_at,omitempty"` // YYYY-MM-DD
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
