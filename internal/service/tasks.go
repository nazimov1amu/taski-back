package service

import (
	"context"
	"log"
	"time"

	"taski_backend/internal/models"
	"taski_backend/internal/repository"
)

type TasksService struct {
	repo *repository.TasksRepository
}

func NewTasksService(repo *repository.TasksRepository) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) Get(ctx context.Context, id string) (models.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting task:", err)
		return models.Task{}, err
	}
	return task, nil
}

func (s *TasksService) GetBulk(ctx context.Context, date time.Time) ([]models.Task, error) {
	tasks, err := s.repo.GetBulk(ctx, date)
	if err != nil {
		log.Println("error getting bulk tasks:", err)
		return nil, err
	}
	return tasks, nil
}

func (s *TasksService) Create(ctx context.Context, task models.Task) (models.Task, error) {
	task, err := s.repo.Create(ctx, task)
	if err != nil {
		log.Println("error creating task:", err)
		return models.Task{}, err
	}
	return task, nil
}

func (s *TasksService) Update(ctx context.Context, task models.Task) (models.Task, error) {
	task, err := s.repo.Update(ctx, task)
	if err != nil {
		log.Println("error updating task:", err)
		return models.Task{}, err
	}
	return task, nil
}

func (s *TasksService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting task:", err)
		return err
	}
	return nil
}

func (s *TasksService) Search(ctx context.Context, query string) ([]models.Task, error) {
	tasks, err := s.repo.Search(ctx, query)
	if err != nil {
		log.Println("error searching tasks:", err)
		return nil, err
	}
	return tasks, nil
}