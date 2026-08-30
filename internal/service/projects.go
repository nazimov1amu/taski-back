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

func (s *ProjectsService) Get(ctx context.Context, id string) (models.Project, error) {
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting project:", err)
		return models.Project{}, err
	}
	return project, nil
}

func (s *ProjectsService) GetBulk(ctx context.Context, limit int, offset int) ([]models.Project, error) {
	projects, err := s.repo.GetBulk(ctx, limit, offset)
	if err != nil {
		log.Println("error getting bulk projects:", err)
		return nil, err
	}
	return projects, nil
}

func (s *ProjectsService) Create(ctx context.Context, project models.Project) (models.Project, error) {
	project, err := s.repo.Create(ctx, project)
	if err != nil {
		log.Println("error creating project:", err)
		return models.Project{}, err
	}
	return project, nil
}

func (s *ProjectsService) Update(ctx context.Context, project models.Project) (models.Project, error) {
	project, err := s.repo.Update(ctx, project)
	if err != nil {
		log.Println("error updating project:", err)
		return models.Project{}, err
	}
	return project, nil
}

func (s *ProjectsService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting project:", err)
		return err
	}
	return nil
}
