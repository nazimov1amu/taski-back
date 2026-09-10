package repository

import (
	"context"
	"database/sql"

	"taski_backend/internal/models"
)

type TasksRepository struct {
	db *sql.DB
}

type TaskFilters struct {
	Date *string `json:"date,omitempty"`
	ProjectID *string `json:"project_id,omitempty"`
	Filter *string `json:"filter,omitempty"`
	Limit *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (r *TasksRepository) Get(ctx context.Context, id string) (models.TaskResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return models.TaskResponse{}, err
	}

	const query = `
		SELECT
			t.id, t.project_id, t.title, t.description, t.completed,
			t.planned_at, t.start_time, t.end_time, p.name
		FROM tasks t
		LEFT JOIN projects p ON p.id = t.project_id AND p.user_id = t.user_id
		WHERE t.id = $1 AND t.user_id = $2
	`
	return scanTaskResponse(r.db.QueryRowContext(ctx, query, id, userID))
}

func (r *TasksRepository) GetBulk(ctx context.Context, filters TaskFilters) ([]models.TaskResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
		SELECT
			t.id, t.project_id, t.title, t.description, t.completed,
			t.planned_at, t.start_time, t.end_time, p.name
		FROM tasks t
		LEFT JOIN projects p ON p.id = t.project_id AND p.user_id = t.user_id
		WHERE t.user_id = $1
		  AND ($2::text IS NULL OR t.planned_at = $2)
		  AND ($3::uuid IS NULL OR t.project_id = $3)
		  AND (
			$4::text IS NULL
			OR t.title ILIKE '%' || $4 || '%'
			OR t.description ILIKE '%' || $4 || '%'
		  )
		ORDER BY t.created_at DESC
		LIMIT $5 OFFSET $6
	`
	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
		filters.Date,
		filters.ProjectID,
		filters.Filter,
		filters.Limit,
		filters.Offset,
	)
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

func (r *TasksRepository) UpsertFromEvent(ctx context.Context, task models.UpsertTaskRequest, tx *sql.Tx) error {
	const query = `
		INSERT INTO tasks (id, user_id, project_id, title, description, planned_at, start_time, end_time, completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			project_id  = EXCLUDED.project_id,
			title       = EXCLUDED.title,
			description = EXCLUDED.description,
			planned_at  = EXCLUDED.planned_at,
			start_time  = EXCLUDED.start_time,
			end_time    = EXCLUDED.end_time,
			completed   = EXCLUDED.completed,
			updated_at  = CURRENT_TIMESTAMP
		WHERE tasks.user_id = EXCLUDED.user_id
	`
	_, err := tx.ExecContext(ctx, query, task.ID, task.UserID, nullStr(task.ProjectID), task.Title, nullStr(task.Description), nullStr(task.PlannedAt), nullStr(task.StartTime), nullStr(task.EndTime), task.Completed)
	if err != nil {
		return err
	}
	return nil
}

func (r *TasksRepository) Delete(ctx context.Context, id string, tx *sql.Tx) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func scanTaskResponse(row scanner) (models.TaskResponse, error) {
	var (
		task        models.TaskResponse
		projectID   sql.NullString
		plannedAt   sql.NullString
		startTime   sql.NullString
		endTime     sql.NullString
		projectName sql.NullString
	)
	if err := row.Scan(
		&task.ID,
		&projectID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&plannedAt,
		&startTime,
		&endTime,
		&projectName,
	); err != nil {
		return models.TaskResponse{}, err
	}

	task.ProjectID = projectID.String
	task.PlannedAt = plannedAt.String
	task.StartTime = startTime.String
	task.EndTime = endTime.String
	task.ProjectName = projectName.String
	return task, nil
}
