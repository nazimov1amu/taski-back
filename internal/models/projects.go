package models

type ProjectFilters struct {
	Name *string `json:"name,omitempty" jsonschema:"Project name or description fragment to search"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" jsonschema:"required,Project name"`
	Description string `json:"description,omitempty" jsonschema:"Project description"`
}

type UpdateProjectRequest struct {
	ID          string `json:"id" jsonschema:"required,Project id"`
	Name        string `json:"name" jsonschema:"required,Project name"`
	Description string `json:"description,omitempty" jsonschema:"Project description"`
}

type ProjectResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Tasks       []TaskResponse `json:"tasks,omitempty"`
}
