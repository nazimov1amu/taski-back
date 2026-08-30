package repository

import (
	"context"
	"database/sql"
	"time"

	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (r *TasksRepository) Get(ctx context.Context, id string) (models.Task, error) {
	row := r.db.QueryRowContext(ctx, db.TaskQueries.Get, id)
	return scanTask(row)
}

func (r *TasksRepository) GetBulk(ctx context.Context, date time.Time) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, db.TaskQueries.GetBulk, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TasksRepository) Create(ctx context.Context, task models.Task) (models.Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		db.TaskQueries.Create,
		nullIfEmpty(task.ProjectID),
		task.Title,
		task.Description,
		task.PlannedAt,
		task.StartTime,
		task.EndTime,
	)
	return scanTask(row)
}

func (r *TasksRepository) Update(ctx context.Context, task models.Task) (models.Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		db.TaskQueries.Update,
		task.ID,
		nullIfEmpty(task.ProjectID),
		task.Title,
		task.Description,
		task.Completed,
		task.PlannedAt,
		task.StartTime,
		task.EndTime,
	)
	return scanTask(row)
}

func (r *TasksRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, db.TaskQueries.Delete, id)
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

func (r *TasksRepository) Search(ctx context.Context, query string) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, db.TaskQueries.Search, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
type scanner interface {
	Scan(dest ...any) error
}

func scanTask(s scanner) (models.Task, error) {
	var (
		task      models.Task
		projectID sql.NullString
		plannedAt time.Time
	)
	err := s.Scan(
		&task.ID,
		&projectID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&plannedAt,
		&task.StartTime,
		&task.EndTime,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return models.Task{}, err
	}
	task.ProjectID = projectID.String
	task.PlannedAt = plannedAt.Format(time.DateOnly)
	return task, nil
}

func nullIfEmpty(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
