package service

import (
	"context"
	"log"

	"taski_backend/internal/models"
	"taski_backend/internal/repository"
)

type TasksService struct {
	repo *repository.TasksRepository
}

func NewTasksService(repo *repository.TasksRepository) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) Get(ctx context.Context, id string) (models.TaskResponse, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting task:", err)
		return models.TaskResponse{}, err
	}
	return task, nil
}

func (s *TasksService) GetBulk(ctx context.Context, args models.TaskFilters) ([]models.TaskResponse, error) {
	tasks, err := s.repo.GetBulk(ctx, args)
	if err != nil {
		log.Println("error getting bulk tasks:", err)
		return nil, err
	}
	return tasks, nil
}

func (s *TasksService) Create(ctx context.Context, task models.CreateTaskRequest) (models.TaskResponse, error) {
	created, err := s.repo.Create(ctx, task)
	if err != nil {
		log.Println("error creating task:", err)
		return models.TaskResponse{}, err
	}
	return created, nil
}

func (s *TasksService) Update(ctx context.Context, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	updated, err := s.repo.Update(ctx, task)
	if err != nil {
		log.Println("error updating task:", err)
		return models.TaskResponse{}, err
	}
	return updated, nil
}

func (s *TasksService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting task:", err)
		return err
	}
	return nil
}

func (s *TasksService) Search(ctx context.Context, query string) ([]models.TaskResponse, error) {
	tasks, err := s.repo.Search(ctx, query)
	if err != nil {
		log.Println("error searching tasks:", err)
		return nil, err
	}
	return tasks, nil
}
