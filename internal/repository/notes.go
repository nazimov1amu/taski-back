package repository

import (
	"context"
	"database/sql"
	"time"

	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type NotesRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

func (r *NotesRepository) Get(ctx context.Context, id string) (models.NoteResponse, error) {
	row := r.db.QueryRowContext(ctx, db.NoteQueries.Get, id)
	return scanNoteResponse(row)
}

func (r *NotesRepository) GetBulk(ctx context.Context) ([]models.NoteResponse, error) {
	rows, err := r.db.QueryContext(ctx, db.NoteQueries.GetBulk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]models.NoteResponse, 0)
	for rows.Next() {
		note, err := scanNoteResponse(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NotesRepository) Create(ctx context.Context, note models.CreateNoteRequest) (models.NoteResponse, error) {
	row := r.db.QueryRowContext(ctx, db.NoteQueries.Create, note.Title, note.Content)
	return scanNoteResponse(row)
}

func (r *NotesRepository) Update(ctx context.Context, id string, note models.UpdateNoteRequest) (models.NoteResponse, error) {
	row := r.db.QueryRowContext(ctx, db.NoteQueries.Update, id, note.Title, note.Content)
	return scanNoteResponse(row)
}

func (r *NotesRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, db.NoteQueries.Delete, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func scanNoteResponse(s scanner) (models.NoteResponse, error) {
	var (
		resp      models.NoteResponse
		createdAt time.Time
		updatedAt time.Time
	)
	err := s.Scan(&resp.ID, &resp.Title, &resp.Content, &createdAt, &updatedAt)
	if err != nil {
		return models.NoteResponse{}, err
	}
	return resp, nil
}
