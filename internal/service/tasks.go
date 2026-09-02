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

func (s *TasksService) Get(ctx context.Context, id string) (models.TaskResponse, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting task:", err)
		return models.TaskResponse{}, mapRepoError(err, "task")
	}
	return task, nil
}

func (s *TasksService) GetBulk(ctx context.Context, args models.TaskFilters) ([]models.TaskResponse, error) {
	tasks, err := s.repo.GetBulk(ctx, args)
	if err != nil {
		log.Println("error getting bulk tasks:", err)
		return nil, mapRepoError(err, "task")
	}
	return tasks, nil
}

func (s *TasksService) Create(ctx context.Context, task models.CreateTaskRequest) (models.TaskResponse, error) {
	created, err := s.repo.Create(ctx, task)
	if err != nil {
		log.Println("error creating task:", err)
		return models.TaskResponse{}, mapRepoError(err, "task")
	}
	return created, nil
}

func (s *TasksService) Update(ctx context.Context, id string, task models.UpdateTaskRequest) (models.TaskResponse, error) {
	if task.StartTime != "" && task.EndTime == "" {
		return models.TaskResponse{}, apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_start_time")
	}

	if task.StartTime != "" && task.EndTime != "" && task.StartTime >= task.EndTime {
		return models.TaskResponse{}, apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_end_time")
	}

	if task.PlannedAt != "" && task.StartTime != "" && task.EndTime != "" {
		tasks, err := s.repo.GetBulk(ctx, models.TaskFilters{
			Date: &task.PlannedAt,
		})
		if err != nil {
			log.Println("error getting tasks:", err)
			return models.TaskResponse{}, mapRepoError(err, "task")
		}
		
		if len(tasks) > 0 {
			for _, t := range tasks {
				if t.StartTime > task.StartTime && t.EndTime < task.EndTime {
					return models.TaskResponse{}, apperrors.NewAppError(apperrors.ErrConflict, "time_range_conflict")
				}
			}
		}
	}

	updated, err := s.repo.Update(ctx, id, task)
	if err != nil {
		log.Println("error updating task:", err)
		return models.TaskResponse{}, mapRepoError(err, "task")
	}
	return updated, nil
}

func (s *TasksService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting task:", err)
		return mapRepoError(err, "task")
	}
	return nil
}

func (s *TasksService) Complete(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(ctx, taskID)
	if err != nil {
		log.Println("error getting task:", err)
		return mapRepoError(err, "task")
	}

	task.TaskFields.Completed = true
	updatedReq := models.UpdateTaskRequest{
		Completed: true,
	}

	_, err = s.repo.Update(ctx, taskID, updatedReq)
	if err != nil {
		log.Println("error updating task:", err)
		return mapRepoError(err, "task")
	}
	return nil
}
