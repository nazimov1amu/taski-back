package repository

import (
	"context"
	"database/sql"

	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type ProjectsRepository struct {
	db *sql.DB
}

func NewProjectsRepository(db *sql.DB) *ProjectsRepository {
	return &ProjectsRepository{db: db}
}

func (r *ProjectsRepository) Get(ctx context.Context, id string) (models.Project, error) {
	var project models.Project
	err := r.db.QueryRowContext(ctx, db.ProjectQueries.Get, id).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return models.Project{}, err
	}
	return project, nil
}

func (r *ProjectsRepository) GetBulk(ctx context.Context, limit int, offset int) ([]models.Project, error) {
	rows, err := r.db.QueryContext(ctx, db.ProjectQueries.GetBulk, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectsRepository) Create(ctx context.Context, project models.Project) (models.Project, error) {
	err := r.db.QueryRowContext(ctx, db.ProjectQueries.Create, project.Name, project.Description).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return models.Project{}, err
	}
	return project, nil
}

func (r *ProjectsRepository) Update(ctx context.Context, project models.Project) (models.Project, error) {
	err := r.db.QueryRowContext(ctx, db.ProjectQueries.Update, project.ID, project.Name, project.Description).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return models.Project{}, err
	}
	return project, nil
}

func (r *ProjectsRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, db.ProjectQueries.Delete, id)
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
