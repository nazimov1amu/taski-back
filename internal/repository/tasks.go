package repository

import (
	"context"
	"database/sql"
	"time"

	"taski_backend/internal/constants"
	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (r *TasksRepository) Get(ctx context.Context, id string) (models.TaskResponse, error) {
	row := r.db.QueryRowContext(ctx, db.TaskQueries.Get, id)
	return scanTaskResponse(row)
}

func (r *TasksRepository) GetBulk(ctx context.Context, args models.TaskFilters) ([]models.TaskResponse, error) {
	rows, err := r.db.QueryContext(ctx, db.TaskQueries.GetBulk, args.Date, args.ProjectID, args.Title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.TaskResponse, 0)
	for rows.Next() {
		task, err := scanTaskResponse(rows)
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

func (r *TasksRepository) Create(ctx context.Context, task models.CreateTaskRequest) (models.TaskResponse, error) {
	row := r.db.QueryRowContext(
		ctx,
		db.TaskQueries.Create,
		nullIfEmpty(task.ProjectID),
		task.Title,
		task.Description,
		task.StartTime,
		task.EndTime,
	)
	return scanTaskResponse(row)
}

func (r *TasksRepository) Update(ctx context.Context, id string, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	row := r.db.QueryRowContext(
		ctx,
		db.TaskQueries.Update,
		id,
		nullIfEmpty(task.ProjectID),
		task.Title,
		task.Description,
		task.Completed,
		task.StartTime,
		task.EndTime,
	)
	return scanTaskResponse(row)
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

type scanner interface {
	Scan(dest ...any) error
}

func scanTaskResponse(s scanner) (models.TaskResponse, error) {
	var (
		resp        models.TaskResponse
		projectID   sql.NullString
		projectName sql.NullString
		startTime   time.Time
		endTime     time.Time
		createdAt   time.Time
		updatedAt   time.Time
	)
	err := s.Scan(
		&resp.ID,
		&projectID,
		&resp.Title,
		&resp.Description,
		&resp.Completed,
		&startTime,
		&endTime,
		&createdAt,
		&updatedAt,
		&projectName,
	)
	if err != nil {
		return models.TaskResponse{}, err
	}
	resp.ProjectID = projectID.String
	resp.ProjectName = projectName.String
	resp.StartTime = startTime.Format(constants.LocalDateTimeLayout)
	resp.EndTime = endTime.Format(constants.LocalDateTimeLayout)
	return resp, nil
}

func nullIfEmpty(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
