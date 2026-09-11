package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"taski_backend/internal/models"
	"taski_backend/internal/service"
)

type CreateEventsResponse struct {
	LastSequenceID int64 `json:"last_sequence_id"`
}
type SyncHandler struct {
	syncService *service.SyncService
}

func NewSyncHandler(syncService *service.SyncService) *SyncHandler {
	return &SyncHandler{syncService: syncService}
}

func (h *SyncHandler) CreateEvents(w http.ResponseWriter, r *http.Request) {
	var events []models.EventCreateRequest
	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		writeServiceError(w, err, "sync")
		return
	}
	defer r.Body.Close()
	lastSequenceID, err := h.syncService.CreateEvents(r.Context(), events)
	if err != nil {
		writeServiceError(w, err, "sync")
		return
	}

	writeJSON(w, http.StatusCreated, CreateEventsResponse{LastSequenceID: lastSequenceID})
}

func (h *SyncHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	sequenceID, err := strconv.ParseInt(r.URL.Query().Get("sequence_id"), 10, 64)
	if err != nil {
		log.Println("error parsing sequence_id:", err)
		writeServiceError(w, err, "sync")
		return
	}
	events, err := h.syncService.GetEvents(r.Context(), sequenceID)
	if err != nil {
		log.Println("error getting events:", err)
		writeServiceError(w, err, "sync")
		return
	}
	writeJSON(w, http.StatusOK, events)
}
