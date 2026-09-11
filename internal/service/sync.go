package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"taski_backend/internal/apperrors"
	"taski_backend/internal/models"
	"taski_backend/internal/repository"

	"go.uber.org/zap"
)

type SyncService struct {
	db                 *sql.DB
	syncRepository     *repository.SyncRepository
	tasksRepository    *repository.TasksRepository
	projectsRepository *repository.ProjectsRepository
	logger             *zap.SugaredLogger
}

func NewSyncService(db *sql.DB, logger *zap.SugaredLogger) *SyncService {
	return &SyncService{
		db:                 db,
		syncRepository:     repository.NewSyncRepository(db),
		tasksRepository:    repository.NewTasksRepository(db),
		projectsRepository: repository.NewProjectsRepository(db),
		logger:             logger,
	}
}

func (s *SyncService) CreateEvents(ctx context.Context, events []models.EventCreateRequest) (int64, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("error beginning transaction", "error", err)
		return 0, err
	}
	defer tx.Rollback()

	sequenceID, err := s.syncRepository.GetCurrentSequenceID(ctx, userID)
	if err != nil {
		s.logger.Error("error getting current sequence ID", "error", err)
		return 0, err
	}
	sequenceID++

	for i := range events {
		events[i].UserID = userID
		events[i].SequenceID = sequenceID
		sequenceID++
	}

	if err := s.syncRepository.CreateEvents(ctx, events, tx); err != nil {
		s.logger.Error("error creating events", "error", err)
		return 0, err
	}

	if err := s.ManageEvents(ctx, events, tx); err != nil {
		s.logger.Error("error managing events", "error", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("error committing transaction", "error", err)
		return 0, err
	}

	lastSequenceID, err := s.syncRepository.GetCurrentSequenceID(ctx, userID)
	if err != nil {
		s.logger.Error("error getting current sequence ID", "error", err)
		return 0, err
	}

	return lastSequenceID, nil
}

func (s *SyncService) GetEvents(ctx context.Context, sequenceID int64) ([]models.Event, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	events, err := s.syncRepository.GetEvents(ctx, sequenceID, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SyncService) ManageEvents(ctx context.Context, events []models.EventCreateRequest, tx *sql.Tx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := s.applyEvent(ctx, userID, event, tx); err != nil {
			s.logger.Error("error applying event", "error", err)
			return err
		}
	}

	return nil
}

func (s *SyncService) applyEvent(ctx context.Context, userID string, event models.EventCreateRequest, tx *sql.Tx) error {
	switch event.EntityType {
	case "task":
		return s.applyTaskEvent(ctx, userID, event, tx)
	case "project":
		return s.applyProjectEvent(ctx, userID, event, tx)
	default:
		return apperrors.NewAppError(apperrors.ErrInvalidInput, "unsupported_entity_type")
	}
}

func (s *SyncService) applyTaskEvent(ctx context.Context, userID string, event models.EventCreateRequest, tx *sql.Tx) error {
	switch event.Operation {
	case "create", "update":
		task, err := decodePayload[models.UpsertTaskRequest](event.Payload, s.logger)
		if err != nil {
			return apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_task_payload")
		}
		task.UserID = userID
		return s.tasksRepository.UpsertFromEvent(ctx, task, tx)
	case "delete":
		return s.tasksRepository.Delete(ctx, event.EntityID, tx)
	default:
		return apperrors.NewAppError(apperrors.ErrInvalidInput, "unsupported_operation")
	}
}

func (s *SyncService) applyProjectEvent(ctx context.Context, userID string, event models.EventCreateRequest, tx *sql.Tx) error {
	switch event.Operation {
	case "create", "update":
		project, err := decodePayload[models.UpsertProjectRequest](event.Payload, s.logger)
		if err != nil {
			return apperrors.NewAppError(apperrors.ErrInvalidInput, "invalid_project_payload")
		}
		project.UserID = userID
		return s.projectsRepository.UpsertFromEvent(ctx, project, tx)
	case "delete":
		return s.projectsRepository.Delete(ctx, event.EntityID, tx)
	default:
		return apperrors.NewAppError(apperrors.ErrInvalidInput, "unsupported_operation")
	}
}

func decodePayload[T any](payload any, logger *zap.SugaredLogger) (T, error) {
	var out T
	if payload == nil {
		return out, fmt.Errorf("empty payload")
	}
	if v, ok := payload.(T); ok {
		return v, nil
	}
	b, err := json.Marshal(payload)
	if err != nil {
		logger.Error("error marshalling payload", "error", err)
		return out, err
	}
	if err := json.Unmarshal(b, &out); err != nil {
		logger.Error("error unmarshalling payload", "error", err)
		return out, err
	}
	return out, nil
}
