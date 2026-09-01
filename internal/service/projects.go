package service

import (
	"context"
	"log"

	"taski_backend/internal/models"
	"taski_backend/internal/repository"
)

type ProjectsService struct {
	repo *repository.ProjectsRepository
}

func NewProjectsService(repo *repository.ProjectsRepository) *ProjectsService {
	return &ProjectsService{repo: repo}
}

func (s *ProjectsService) Get(ctx context.Context, id string) (models.ProjectResponse, error) {
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting project:", err)
		return models.ProjectResponse{}, mapRepoError(err)
	}
	return project, nil
}

func (s *ProjectsService) GetBulk(ctx context.Context, args models.ProjectFilters) ([]models.ProjectResponse, error) {
	projects, err := s.repo.GetBulk(ctx, args)
	if err != nil {
		log.Println("error getting bulk projects:", err)
		return nil, mapRepoError(err)
	}
	return projects, nil
}

func (s *ProjectsService) Create(ctx context.Context, project models.CreateProjectRequest) (models.ProjectResponse, error) {
	created, err := s.repo.Create(ctx, project)
	if err != nil {
		log.Println("error creating project:", err)
		return models.ProjectResponse{}, mapRepoError(err)
	}
	return created, nil
}

func (s *ProjectsService) Update(ctx context.Context, id string, project models.UpdateProjectRequest) (models.ProjectResponse, error) {
	updated, err := s.repo.Update(ctx, id, project)
	if err != nil {
		log.Println("error updating project:", err)
		return models.ProjectResponse{}, mapRepoError(err)
	}
	return updated, nil
}

func (s *ProjectsService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting project:", err)
		return mapRepoError(err)
	}
	return nil
}
