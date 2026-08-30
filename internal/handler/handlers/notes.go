package handlers

import (
	"encoding/json"
	"net/http"

	"taski_backend/internal/handler/dto"
	"taski_backend/internal/models"
	"taski_backend/internal/service"
)

type NotesHandler struct {
	service *service.NotesService
}

func NewNotesHandler(service *service.NotesService) *NotesHandler {
	return &NotesHandler{service: service}
}

func (h *NotesHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	note, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toNoteResponse(note))
}

func (h *NotesHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notes, err := h.service.GetBulk(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.NoteResponse, 0, len(notes))
	for _, note := range notes {
		resp = append(resp, toNoteResponse(note))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(r.Context(), models.Note{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toNoteResponse(created))
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := parseUUID(req.ID); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(r.Context(), models.Note{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toNoteResponse(updated))
}

func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toNoteResponse(note models.Note) dto.NoteResponse {
	return dto.NoteResponse{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}
}
