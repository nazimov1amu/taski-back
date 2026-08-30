package service

import (
	"context"
	"log"

	"taski_backend/internal/models"
	"taski_backend/internal/repository"
)

type NotesService struct {
	repo *repository.NotesRepository
}

func NewNotesService(repo *repository.NotesRepository) *NotesService {
	return &NotesService{repo: repo}
}

func (s *NotesService) Get(ctx context.Context, id string) (models.Note, error) {
	note, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting note:", err)
		return models.Note{}, err
	}
	return note, nil
}

func (s *NotesService) GetBulk(ctx context.Context, limit int, offset int) ([]models.Note, error) {
	notes, err := s.repo.GetBulk(ctx, limit, offset)
	if err != nil {
		log.Println("error getting bulk notes:", err)
		return nil, err
	}
	return notes, nil
}

func (s *NotesService) Create(ctx context.Context, note models.Note) (models.Note, error) {
	note, err := s.repo.Create(ctx, note)
	if err != nil {
		log.Println("error creating note:", err)
		return models.Note{}, err
	}
	return note, nil
}

func (s *NotesService) Update(ctx context.Context, note models.Note) (models.Note, error) {
	note, err := s.repo.Update(ctx, note)
	if err != nil {
		log.Println("error updating note:", err)
		return models.Note{}, err
	}
	return note, nil
}

func (s *NotesService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting note:", err)
		return err
	}
	return nil
}
