package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"taski_backend/internal/models"
)

type SyncRepository struct {
	db *sql.DB
}

func NewSyncRepository(db *sql.DB) *SyncRepository {
	return &SyncRepository{db: db}
}

func (r *SyncRepository) CreateEvents(ctx context.Context, events []models.EventCreateRequest, tx *sql.Tx) error {
	var q strings.Builder
	var args []any

	q.WriteString("INSERT INTO sync_events (user_id, sequence_id, entity_type, entity_id, operation, payload) VALUES ")

	for i, e := range events {
		if i > 0 {
			q.WriteByte(',')
		}
		n := i * 8
		fmt.Fprintf(&q, "($%d,$%d,$%d,$%d,$%d,$%d)",
			n+1, n+2, n+3, n+4, n+5, n+6)
		args = append(args,
			e.UserID, e.SequenceID, e.EntityType, e.EntityID, e.Operation, e.Payload,
		)
	}

	query := q.String()
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error creating events:", err)
		return err
	}
	return nil
}

func (r *SyncRepository) GetEvents(ctx context.Context, sequenceID int64, userID string) ([]models.Event, error) {
	query := `
		SELECT id, user_id, sequence_id, entity_type, entity_id, operation, payload, created_at FROM sync_events
		WHERE sequence_id > $1
		AND user_id = $2
		ORDER BY sequence_id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, sequenceID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *SyncRepository) GetCurrentSequenceID(ctx context.Context, userID string) (int64, error) {
	query := `
		SELECT sequence_id FROM sync_events
		WHERE user_id = $1
		ORDER BY sequence_id DESC
		LIMIT 1
	`
	var sequenceID int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sequenceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return sequenceID, nil
}

func scanEvent(row scanner) (models.Event, error) {
	var event models.Event
	err := row.Scan(&event.ID, &event.UserID, &event.SequenceID, &event.EntityType, &event.EntityID, &event.Operation, &event.Payload, &event.CreatedAt)
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}