package models


type UpsertTaskRequest struct {
	ID          string `json:"id"`
	UserID string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	PlannedAt   string `json:"planned_at,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	ProjectName string `json:"project_name,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	PlannedAt   string `json:"planned_at,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
}
