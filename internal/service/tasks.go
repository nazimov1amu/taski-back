package service

import (
	"context"
	"log"

	"taski_backend/internal/apperrors"
	"taski_backend/internal/models"
	"taski_backend/internal/repository"
)

type TasksService struct {
	repo *repository.TasksRepository
}

func NewTasksService(repo *repository.TasksRepository) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) Get(ctx context.Context, id string) (models.TaskWithChildrenResponse, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting task:", err)
		return models.TaskWithChildrenResponse{}, mapRepoError(err)
	}
	return task, nil
}

func (s *TasksService) GetBulk(ctx context.Context, args models.TaskFilters) ([]models.TaskResponse, error) {
	tasks, err := s.repo.GetBulk(ctx, args)
	if err != nil {
		log.Println("error getting bulk tasks:", err)
		return nil, mapRepoError(err)
	}
	return tasks, nil
}

func (s *TasksService) Create(ctx context.Context, task models.CreateTaskRequest) (models.TaskResponse, error) {
	created, err := s.repo.Create(ctx, task)
	if err != nil {
		log.Println("error creating task:", err)
		return models.TaskResponse{}, mapRepoError(err)
	}
	return created, nil
}

func (s *TasksService) Update(ctx context.Context, id string, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	if task.ParentID == id {
		return models.TaskResponse{}, apperrors.ErrInvalidInput
	}

	updated, err := s.repo.Update(ctx, id, task)
	if err != nil {
		log.Println("error updating task:", err)
		return models.TaskResponse{}, mapRepoError(err)
	}
	return updated, nil
}

func (s *TasksService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting task:", err)
		return mapRepoError(err)
	}
	return nil
}

func (s *TasksService) Complete(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(ctx, taskID)
	if err != nil {
		log.Println("error getting task:", err)
		return mapRepoError(err)
	}

	task.TaskFields.Completed = true
	updatedReq := models.UpdateTaskRequest{
		TaskFields: task.TaskFields,
	}

	_, err = s.repo.Update(ctx, taskID, updatedReq)
	if err != nil {
		log.Println("error updating task:", err)
		return mapRepoError(err)
	}
	return nil
}
