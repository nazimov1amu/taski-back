package repository

import (
	"context"
	"database/sql"

	"taski_backend/internal/models"
)

type ProjectsRepository struct {
	db *sql.DB
}

type ProjectFilters struct {
	Name *string `json:"name,omitempty"`
}

func NewProjectsRepository(db *sql.DB) *ProjectsRepository {
	return &ProjectsRepository{db: db}
}


func (r *ProjectsRepository) Get(ctx context.Context, id string) (models.ProjectResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return models.ProjectResponse{}, err
	}

	const query = `
		SELECT id, name, description
		FROM projects
		WHERE id = $1 AND user_id = $2
	`
	return scanProjectResponse(r.db.QueryRowContext(ctx, query, id, userID))
}

func (r *ProjectsRepository) GetBulk(ctx context.Context, filters ProjectFilters) ([]models.ProjectResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
		SELECT id, name, description
		FROM projects
		WHERE user_id = $1
		  AND (
			$2::text IS NULL
			OR name ILIKE '%' || $2 || '%'
			OR description ILIKE '%' || $2 || '%'
		  )
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, filters.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.ProjectResponse, 0)
	for rows.Next() {
		project, err := scanProjectResponse(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectsRepository) UpsertFromEvent(ctx context.Context, project models.UpsertProjectRequest, tx *sql.Tx) error {
	const query = `
		INSERT INTO projects (id, user_id, name, description)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name        = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at  = CURRENT_TIMESTAMP
		WHERE projects.user_id = EXCLUDED.user_id
	`
	_, err := tx.ExecContext(ctx, query, project.ID, project.UserID, project.Name, nullStr(project.Description))
	return err
}

func (r *ProjectsRepository) Delete(ctx context.Context, id string, tx *sql.Tx) error {
	const query = `DELETE FROM projects WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func scanProjectResponse(row scanner) (models.ProjectResponse, error) {
	var project models.ProjectResponse
	if err := row.Scan(&project.ID, &project.Name, &project.Description); err != nil {
		return models.ProjectResponse{}, err
	}
	return project, nil
}

