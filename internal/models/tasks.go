package models

type TaskFilters struct {
	Date      *string
	ProjectID *string
	Title     *string
}

type CreateTaskRequest struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectID   string    `json:"project_id,omitempty"`
	StartTime   string    `json:"start_time,omitempty"`
	EndTime     string    `json:"end_time,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectID   string    `json:"project_id,omitempty"`
	Completed   bool      `json:"completed"`
	StartTime   string    `json:"start_time,omitempty"`
	EndTime     string    `json:"end_time,omitempty"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ProjectID   string    `json:"project_id,omitempty"`
	ProjectName string    `json:"project_name,omitempty"`
	Completed   bool      `json:"completed"`
	StartTime   string    `json:"start_time,omitempty"`
	EndTime     string    `json:"end_time,omitempty"`
}
