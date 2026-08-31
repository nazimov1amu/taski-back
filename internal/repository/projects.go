package repository

import (
	"context"
	"database/sql"
	"time"

	db "taski_backend/internal/db/queries"
	"taski_backend/internal/models"
)

type ProjectsRepository struct {
	db *sql.DB
}

func NewProjectsRepository(db *sql.DB) *ProjectsRepository {
	return &ProjectsRepository{db: db}
}

func (r *ProjectsRepository) Get(ctx context.Context, id string) (models.ProjectResponse, error) {
	row := r.db.QueryRowContext(ctx, db.ProjectQueries.Get, id)
	return scanProjectResponse(row)
}

func (r *ProjectsRepository) GetBulk(ctx context.Context) ([]models.ProjectResponse, error) {
	rows, err := r.db.QueryContext(ctx, db.ProjectQueries.GetBulk)
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

func (r *ProjectsRepository) Create(ctx context.Context, project models.CreateProjectRequest) (models.ProjectResponse, error) {
	row := r.db.QueryRowContext(ctx, db.ProjectQueries.Create, project.Name, project.Description)
	return scanProjectResponse(row)
}

func (r *ProjectsRepository) Update(ctx context.Context, id string, project models.UpdateProjectRequest) (models.ProjectResponse, error) {
	row := r.db.QueryRowContext(ctx, db.ProjectQueries.Update, id, project.Name, project.Description)
	return scanProjectResponse(row)
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

func scanProjectResponse(s scanner) (models.ProjectResponse, error) {
	var (
		resp      models.ProjectResponse
		createdAt time.Time
		updatedAt time.Time
	)
	err := s.Scan(&resp.ID, &resp.Name, &resp.Description, &createdAt, &updatedAt)
	if err != nil {
		return models.ProjectResponse{}, err
	}
	return resp, nil
}
