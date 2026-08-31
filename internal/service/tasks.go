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

func (s *TasksService) Update(ctx context.Context, id string, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	updated, err := s.repo.Update(ctx, id, task)
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


func (s *TasksService) Complete(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(ctx, taskID)
	if err != nil {
		log.Println("error getting task:", err)
		return err
	}
	
	updatedReq := models.UpdateTaskRequest{
		Title:     task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
		PlannedAt:   task.PlannedAt,
		StartTime:   task.StartTime,
		EndTime:     task.EndTime,
		Completed: true,
	}

	_, err = s.repo.Update(ctx, taskID, updatedReq)
	if err != nil {
		log.Println("error updating task:", err)
		return err
	}
	return nil
}
