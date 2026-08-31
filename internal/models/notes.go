package models

type CreateNoteRequest struct {
	Title   string `json:"title" jsonschema:"required,Note title"`
	Content string `json:"content,omitempty" jsonschema:"Note content"`
}

type UpdateNoteRequest struct {
	ID      string `json:"id" jsonschema:"required,Note id"`
	Title   string `json:"title" jsonschema:"required,Note title"`
	Content string `json:"content,omitempty" jsonschema:"Note content"`
}

type NoteResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
