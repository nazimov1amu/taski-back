package models

type TaskFilters struct {
	Date      *string
	ProjectID *string
	Title     *string
}

type TaskFields struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	PlannedAt   string `json:"planned_at,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
	Completed   bool   `json:"completed"`
}

type CreateTaskRequest struct {
	TaskFields
}

type UpdateTaskRequest struct {
	TaskFields
}

type TaskResponse struct {
	ID          string `json:"id"`
	ProjectName string `json:"project_name,omitempty"`
	TaskFields
}
