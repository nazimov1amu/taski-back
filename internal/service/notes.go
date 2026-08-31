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

func (s *NotesService) Get(ctx context.Context, id string) (models.NoteResponse, error) {
	note, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Println("error getting note:", err)
		return models.NoteResponse{}, err
	}
	return note, nil
}

func (s *NotesService) GetBulk(ctx context.Context) ([]models.NoteResponse, error) {
	notes, err := s.repo.GetBulk(ctx)
	if err != nil {
		log.Println("error getting bulk notes:", err)
		return nil, err
	}
	return notes, nil
}

func (s *NotesService) Create(ctx context.Context, note models.CreateNoteRequest) (models.NoteResponse, error) {
	created, err := s.repo.Create(ctx, note)
	if err != nil {
		log.Println("error creating note:", err)
		return models.NoteResponse{}, err
	}
	return created, nil
}

func (s *NotesService) Update(ctx context.Context, note models.UpdateNoteRequest) (models.NoteResponse, error) {
	updated, err := s.repo.Update(ctx, note)
	if err != nil {
		log.Println("error updating note:", err)
		return models.NoteResponse{}, err
	}
	return updated, nil
}

func (s *NotesService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Println("error deleting note:", err)
		return err
	}
	return nil
}
