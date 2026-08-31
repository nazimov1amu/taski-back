package handlers

import (
	"encoding/json"
	"net/http"

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

	writeJSON(w, http.StatusOK, note)
}

func (h *NotesHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetBulk(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := parseUUID(req.ID); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
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
