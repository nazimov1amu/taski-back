package models

type ProjectFilters struct {
	Name *string `json:"name,omitempty" jsonschema:"Project name or description fragment to search"`
}

type TaskFilters struct {
	Date      *string `json:"date,omitempty"`
	ProjectID *string `json:"project_id,omitempty"`
	Filter    *string `json:"filter,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
	Offset    *int    `json:"offset,omitempty"`
}
