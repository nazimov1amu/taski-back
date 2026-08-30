package repository

import (
	"context"
	"database/sql"

	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type NotesRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

func (r *NotesRepository) Get(ctx context.Context, id string) (models.Note, error) {
	var note models.Note
	err := r.db.QueryRowContext(ctx, db.NoteQueries.Get, id).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func (r *NotesRepository) GetBulk(ctx context.Context, limit int, offset int) ([]models.Note, error) {
	rows, err := r.db.QueryContext(ctx, db.NoteQueries.GetBulk, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NotesRepository) Create(ctx context.Context, note models.Note) (models.Note, error) {
	err := r.db.QueryRowContext(ctx, db.NoteQueries.Create, note.Title, note.Content).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func (r *NotesRepository) Update(ctx context.Context, note models.Note) (models.Note, error) {
	err := r.db.QueryRowContext(ctx, db.NoteQueries.Update, note.ID, note.Title, note.Content).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return models.Note{}, err
	}
	return note, nil
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
