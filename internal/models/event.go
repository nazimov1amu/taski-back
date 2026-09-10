package models

type EventCreateRequest struct {
	UserID      string `json:"user_id,omitempty"`
	SequenceID  int64 `json:"sequence_id,omitempty"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Operation   string `json:"operation"`
	Payload     any `json:"payload,omitempty"`
}


type Event struct {
	ID          int64 `json:"id"`
	UserID      string `json:"user_id"`
	SequenceID  int64 `json:"sequence_id"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Operation   string `json:"operation"`
	Payload     any `json:"payload"`
	CreatedAt   string `json:"created_at"`
}