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
		nullIfEmpty(task.Description),
		nullIfEmpty(task.PlannedAt),
		nullIfEmpty(task.StartTime),
		nullIfEmpty(task.EndTime),
	)
	return scanTaskResponse(row)
}

func (r *TasksRepository) Update(ctx context.Context, id string, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	row := r.db.QueryRowContext(
		ctx,
		db.TaskQueries.Update,
		id,
		nullIfEmpty(task.ProjectID),
		nullIfEmpty(task.Title),
		nullIfEmpty(task.Description),
		nullIfEmpty(task.Completed),
		nullIfEmpty(task.PlannedAt),
		nullIfEmpty(task.StartTime),
		nullIfEmpty(task.EndTime),
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
		description sql.NullString
		plannedAt   sql.NullString
		startTime   sql.NullTime
		endTime     sql.NullTime
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := s.Scan(
		&resp.ID,
		&projectID,
		&resp.Title,
		&description,
		&resp.Completed,
		&plannedAt,
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
	resp.Description = description.String
	resp.PlannedAt = plannedAt.String
	if startTime.Valid {
		resp.StartTime = startTime.Time.Format(constants.LocalDateTimeLayout)
	}
	if endTime.Valid {
		resp.EndTime = endTime.Time.Format(constants.LocalDateTimeLayout)
	}
	return resp, nil
}

func nullIfEmpty(s any) any{
	switch v := s.(type) {
	case string:
		if v == "" {
			return sql.NullString{}
		}
		return sql.NullString{String: v, Valid: true}
	case bool:
		if !v {
			return sql.NullBool{}
		}
		return sql.NullBool{Bool: v, Valid: true}
	case time.Time:
		if v.IsZero() {
			return sql.NullTime{}
		}
		return sql.NullTime{Time: v, Valid: true}
	case nil:
		return sql.NullString{}
	}
	return s
}
